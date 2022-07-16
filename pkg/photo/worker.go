package photo

import (
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
			reader, err := os.Open(path.Join(p.Path, p.FileName))
			if err != nil {
				log.Errorf(err.Error())
			}
			writer, err := os.Create(path.Join(p.GetTargetPath(), p.FileName))
			if err != nil {
				log.Errorf(err.Error())
			}
			bytesWritten, err := io.Copy(writer, reader)
			if err != nil {
				log.Errorf("error copying file %v. err=%v", p.FileName, err.Error())
			}
			writer.Sync()
			reader.Close()
			writer.Close()

			os.Chtimes(writer.Name(), p.Atime, p.Mtime)

			w.Copier.incrementStats(bytesWritten)
			w.Copier.ProgressBar.Add(1)
		}
	}
}
