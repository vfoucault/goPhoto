package photo

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Photo struct {
	Path      string
	FileName  string
	DateTaken time.Time
	Copier    *Copier
	Atime     time.Time
	Ctime     time.Time
	Mtime     time.Time
	Md5       []byte
	File      *os.File
}

func (p *Photo) GetTargetPath() string {
	return path.Join(p.Copier.Config.DestDirectory, p.DateTaken.Format(p.Copier.Config.DestFileFormat))
}

func (p *Photo) Open() error {
	if p.File == nil {
		f, err := os.Open(path.Join(p.Path, p.FileName))
		if err != nil {
			return fmt.Errorf("unable to open file %v. err=%v", p.Path, err.Error())
		}
		p.File = f
	}
	// rewind the file
	p.File.Seek(0, 0)
	return nil
}

func (p *Photo) GetHash() error {
	if err := p.Open(); err != nil {
		return err
	}
	h := md5.New()
	if _, err := io.Copy(h, p.File); err != nil {
		return fmt.Errorf("unable to compute md5 for file %s. err=%v", path.Join(p.Path, p.FileName), err.Error())
	}
	p.Md5 = h.Sum(nil)
	return nil
}

func (p *Photo) GetDateTaken() error {
	if err := p.Open(); err != nil {
		return err
	}
	exifData, err := exif.Decode(p.File)
	if err != nil {
		return fmt.Errorf("unable to decode exif for file %v. err=%v", p.Path, err.Error())
	}

	p.DateTaken, _ = exifData.DateTime()

	return nil
}
