package photo

import (
	"fmt"
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
}

func (p *Photo) GetTargetPath() string {
	return path.Join(p.Copier.Config.DestDirectory, p.DateTaken.Format(p.Copier.Config.DestFileFormat))
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
