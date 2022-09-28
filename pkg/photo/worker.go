package photo

import (
	"crypto/md5"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

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
	log.Debugf("Starting worker id=%v", w.ID)
	for {
		select {
		case <-w.Copier.Context.Done():
			log.Debugf("Stopping worker id=%v", w.ID)
			w.Copier.Wg.Done()
			return nil
		case p := <-w.Copier.CopyQueue:
			log.Debugf("copying file %v to %v", p.FileName, p.GetTargetPath())
			//Check if a file already exists at destination
			if ok := w.CheckSameContents(p); !ok {
				w.Copy(p)
			} else {
				w.Copier.IncrementSkipped()
				w.Copier.ProgressBar.Add(1)
			}
		}
	}
}

func (w *Worker) Copy(p *Photo) {
	// rewind the file
	p.File.Seek(0, 0)
	writer, err := os.Create(path.Join(p.GetTargetPath(), p.FileName))
	if err != nil {
		log.Errorf(err.Error())
	}
	bytesWritten, err := io.Copy(writer, p.File)
	if err != nil {
		log.Errorf("error copying file %v. err=%v", p.FileName, err.Error())
	}
	writer.Sync()
	p.File.Close()
	writer.Close()

	os.Chtimes(writer.Name(), p.Atime, p.Mtime)

	w.Copier.IncrementStats(bytesWritten)
	w.Copier.ProgressBar.Add(1)
}

func (w *Worker) CheckSameContents(p *Photo) bool {
	if _, err := os.Stat(path.Join(p.GetTargetPath(), p.FileName)); err == nil {
		// read and compute md5
		f, err := os.Open(path.Join(p.GetTargetPath(), p.FileName))
		defer f.Close()
		if err != nil {
			return false
		} else {
			h := md5.New()
			if _, err := io.Copy(h, f); err != nil {
				return false
			}
			sum := h.Sum(nil)
			if string(sum) != string(p.Md5) {
				return false
			} else {
				return true
			}
		}
	} else {
		return false
	}

}
