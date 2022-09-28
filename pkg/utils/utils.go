package utils

import (
	"image"
	"os"
	"path"
	"regexp"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
)

var (
	regexImage, _ = regexp.Compile(`^.+(?i)(jpe?g|gif|png|tiff)$`)
)

func IsImage(file os.FileInfo) bool {
	if !file.IsDir() {
		if regexImage.MatchString(file.Name()) {
			return true
		}
		return false
	}
	return false
}

func SaveImage(filePath, fileName string, img image.Image) error {
	log.Infof("saving image %s", path.Join(filePath, fileName))
	return gg.SaveJPG(path.Join(filePath, fileName), img, 90)
}
