package utils

import (
	"os"
	"regexp"
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
