package photo

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/vfoucault/goPhoto/pkg/config"
	"github.com/vfoucault/goPhoto/pkg/utils"
)

type Copier struct {
	Config *config.Config
	Photos []*Photo
	Stats  struct {
		Count   int
		Skipped int
		Size    int64
	}
	StatsMutex  sync.Mutex
	Workers     []*Worker
	CopyQueue   chan *Photo
	Context     context.Context
	CancelFunc  context.CancelFunc
	Wg          sync.WaitGroup
	ProgressBar *progressbar.ProgressBar
}

func (c *Copier) incrementStats(size int64) {
	c.StatsMutex.Lock()
	defer c.StatsMutex.Unlock()
	c.Stats.Count += 1
	c.Stats.Size += size
}

func (c *Copier) incrementSkipped() {
	c.StatsMutex.Lock()
	defer c.StatsMutex.Unlock()
	c.Stats.Skipped += 1
}

func (c *Copier) CreateDestDirs() {
	var dirs = make(map[string]int)
	for _, x := range c.Photos {
		dirs[x.GetTargetPath()] += 1
	}
	for k, _ := range dirs {
		log.Debugf("creating directory %v", k)
		err := os.MkdirAll(k, 0750)
		if err != nil {
			log.Errorf("unable to create directory %v. err=%v", k, err.Error())
		}
	}
}

func NewCopier(config *config.Config, pctx context.Context) *Copier {
	ctx, cancel := context.WithCancel(pctx)
	return &Copier{
		Config:     config,
		CopyQueue:  make(chan *Photo, 10000),
		Context:    ctx,
		CancelFunc: cancel,
	}
}

func (c *Copier) Wait() {
	for {
		if len(c.Photos) == c.Stats.Count+c.Stats.Skipped {
			c.Stop()
			break
		}
		time.Sleep(time.Second)
	}
	c.Wg.Wait()
}

func (c *Copier) Start() error {
	// check if target directory exists
	_, err := os.Stat(c.Config.DestDirectory)

	if err != nil {
		if err == os.ErrNotExist {
			if err = os.MkdirAll(c.Config.DestDirectory, 0750); err != nil {
				return fmt.Errorf("unable to create target directory %v. err=%v", c.Config.DestDirectory, err.Error())
			}
		}
	}

	for _, x := range c.Photos {
		c.CopyQueue <- x
	}

	// launch workers
	for i := 0; i < c.Config.Workers; i++ {
		worker := NewWorker(i, c)
		c.Wg.Add(1)
		c.Workers = append(c.Workers, worker)
		go worker.Start()
	}
	return nil
}

func (c *Copier) Stop() {
	log.Debugf("Stopping copier")
	c.CancelFunc()
}

func (c *Copier) Search() {
	if c.Config.NoRecurse {
		files, _ := ioutil.ReadDir(c.Config.SourceDirectory)
		for _, f := range files {
			if utils.IsImage(f) {
				c.addPhoto(f, c.Config.SourceDirectory)
			}
		}
	} else {
		filepath.Walk(c.Config.SourceDirectory, func(aPath string, f os.FileInfo, _ error) error {
			if utils.IsImage(f) {
				c.addPhoto(f, strings.TrimSuffix(aPath, f.Name()))
			}
			return nil
		})
	}
}

func (c *Copier) addPhoto(f fs.FileInfo, fPath string) {
	photo := &Photo{Path: fPath, FileName: f.Name(), Copier: c}
	stat := f.Sys().(*syscall.Stat_t)
	photo.Atime = time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
	photo.Ctime = time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)
	photo.Mtime = time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
	if err := photo.GetDateTaken(); err != nil {
		log.Errorf("unable to get image date for image %v. err=%v", path.Join(c.Config.SourceDirectory, f.Name()), err.Error())
	} else {
		err := photo.GetHash()
		if err != nil {
			log.Errorf("unable to get hash for file %s", f.Name())
		}
		c.Photos = append(c.Photos, photo)
	}
}

func RunCopier(cfg *config.Config) {
	start := time.Now()

	cfg.PrintConfig()

	ctx := context.Background()
	copier := NewCopier(cfg, ctx)

	copier.Search()
	copier.CreateDestDirs()
	copier.InitProgressBar()

	// Handle stop and more
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go handleSignals(signalChannel, copier)

	log.Debugf("Will have to copy %v pictures", len(copier.Photos))

	if err := copier.Start(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
	copier.Wait()

	elapsed := time.Since(start)
	fmt.Println()
	log.Infof("Copy ended. Took %v", elapsed)
	log.Infof("Copied %d images / %s.", copier.Stats.Count, bytefmt.ByteSize(uint64(copier.Stats.Size)))
	if copier.Stats.Skipped > 0 {
		log.Infof("Skipped %d images that where already present", copier.Stats.Skipped)
	}
	log.Infof("Byte rate %v/s", bytefmt.ByteSize(uint64(copier.Stats.Size/int64(elapsed.Seconds()))))
}

func (c *Copier) InitProgressBar() {
	c.ProgressBar = progressbar.New(len(c.Photos))
}

//func handleSignals(signChannel chan os.Signal, processorManager *models.ProcessorManager) {
func handleSignals(signChannel chan os.Signal, copier *Copier) {
	log.Debugf("Running signal handler")
	// Handle stop and more
	for {
		sig := <-signChannel
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			log.Infof("Received signal %v", sig)
			log.Infof("Shutting down application. Kill running/pending jobs")
			copier.Stop()
			return
		default:
			log.Errorf("Unable to handle signal %v", sig.String())
		}
	}
}
