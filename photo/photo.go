package photo

import (
	"context"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/schollz/progressbar/v3"
	"github.com/vfoucault/goPhoto/config"
	"github.com/vfoucault/goPhoto/logger"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var (
	regexImage, _ = regexp.Compile(`^.+(?i)(jpe?g|gif|png|tiff)$`)
)

type Copier struct {
	Config      *config.Config
	Photos      []*Photo
	Stats       int
	StatsMutex  sync.Mutex
	Workers     []*Worker
	CopyQueue   chan *Photo
	Context     context.Context
	CancelFunc  context.CancelFunc
	Wg          sync.WaitGroup
	ProgressBar *progressbar.ProgressBar
}

func (c *Copier) IncrementStats() {
	c.StatsMutex.Lock()
	c.Stats += 1
	c.StatsMutex.Unlock()
}

type Photo struct {
	Path      string
	FileName  string
	DateTaken time.Time
	Copier    *Copier
	Atime     time.Time
	Ctime     time.Time
	Mtime     time.Time
}

func (p *Photo) GetTargetPath() string {
	return path.Join(p.Copier.Config.DestDirectory, p.DateTaken.Format(p.Copier.Config.DestFileFormat))
}

func NewCopier(config *config.Config, ctx context.Context) *Copier {
	copier := new(Copier)
	copier.Config = config
	copier.CopyQueue = make(chan *Photo, 10000)

	// Search for images
	copier.Search()

	// Create cancel func
	copier.Context, copier.CancelFunc = context.WithCancel(ctx)

	// Create required directories

	copier.ProgressBar = progressbar.New(len(copier.Photos))

	var dirs = make(map[string]int)
	for _, x := range copier.Photos {
		dirs[x.GetTargetPath()] += 1
	}
	for k, _ := range dirs {
		logger.Debugf("creating directory %v", k)
		err := os.MkdirAll(k, 0750)
		if err != nil {
			logger.Errorf("unable to create directory %v. err=%v", k, err.Error())
		}
	}
	return copier
}

func isImage(file os.FileInfo) bool {
	if !file.IsDir() {
		if regexImage.MatchString(file.Name()) {
			return true
		}
		return false
	}
	return false
}

func (c *Copier) Search() {
	if c.Config.NoRecurse {
		files, _ := ioutil.ReadDir(c.Config.SourceDirectory)
		for _, f := range files {
			if isImage(f) {
				photo := &Photo{Path: c.Config.SourceDirectory, FileName: f.Name(), Copier: c}
				//stat := f.Sys().(*syscall.Stat_t)
				//photo.Atime = time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
				//photo.Ctime = time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
				//photo.Mtime = time.Unix(int64(stat.Mtimespec.Sec), int64(stat.Mtimespec.Nsec))
				if err := photo.GetDateTaken(); err != nil {
					logger.Errorf("unable to get image date for image %v. err=%v", path.Join(c.Config.SourceDirectory, f.Name()), err.Error())
				} else {
					c.Photos = append(c.Photos, photo)
				}
			}
		}
	} else {
		filepath.Walk(c.Config.SourceDirectory, func(aPath string, f os.FileInfo, _ error) error {
			if isImage(f) {
				photo := &Photo{Path: filepath.Dir(aPath), FileName: f.Name(), Copier: c}
				if err := photo.GetDateTaken(); err != nil {
					logger.Errorf("unable to get image date for image %v. err=%v", path.Join(c.Config.SourceDirectory, f.Name()), err.Error())
				} else {
					logger.Debugf("adding photo %v", f.Name())
					c.Photos = append(c.Photos, photo)
				}
			}
			//logger.Infof(path.Join(c.Config.SourceDirectory, f.Name()))
			return nil
		})
	}
}

func (p *Photo) GetDateTaken() error {
	f, err := os.Open(p.Path + "/" + p.FileName)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("unable to open file %v. err=%v", p.Path, err.Error())
	}

	exifData, err := exif.Decode(f)
	if err != nil {
		return fmt.Errorf("unable to decode exif for file %v. err=%v", p.Path, err.Error())
	}

	p.DateTaken, _ = exifData.DateTime()
	return nil
}

func (c *Copier) Wait() {
	for {
		if len(c.Photos) == c.Stats {
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
	logger.Debugf("Stopping copier")
	c.CancelFunc()
}

func NewWorker(id int, copier *Copier) *Worker {
	return &Worker{
		ID:     id,
		Copier: copier,
	}
}

type Worker struct {
	ID     int
	Copier *Copier
}

func (w *Worker) Start() error {
	logger.Debugf("Starting worker id=%v", w.ID)
	for {
		select {
		case <-w.Copier.Context.Done():
			logger.Debugf("Stopping worker id=%v", w.ID)
			w.Copier.Wg.Done()
			return nil
		case p := <-w.Copier.CopyQueue:
			logger.Debugf("copying file %v to %v", p.FileName, p.GetTargetPath())
			reader, err := os.Open(path.Join(p.Path, p.FileName))
			if err != nil {
				logger.Errorf(err.Error())
			}
			writer, err := os.Create(path.Join(p.GetTargetPath(), p.FileName))
			if err != nil {
				logger.Errorf(err.Error())
			}
			_, err = io.Copy(writer, reader)
			if err != nil {
				logger.Errorf("error copying file %v. err=%v", p.FileName, err.Error())
			}
			writer.Sync()
			reader.Close()
			writer.Close()

			w.Copier.IncrementStats()
			w.Copier.ProgressBar.Add(1)
		}
	}
}
