package resize

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flopp/go-findfont"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"github.com/vfoucault/goPhoto/pkg/utils"
)

type WaterMark struct {
	Color color.Gray16
	Size  float64
	Text  string
}

func PhotoResize(w, h uint, srcPath, dstPath string, addText bool, wm ...WaterMark) error {
	// list all images
	var hasErrors bool
	filepath.Walk(srcPath, func(aPath string, f os.FileInfo, _ error) error {
		if utils.IsImage(f) {
			log.Infof("resizing %s...", f.Name())
			imagePath := path.Join(srcPath, f.Name())
			img, err := gg.LoadImage(imagePath)
			if err != nil {
				log.Errorf("unable to load image %s. err=%w", imagePath, err.Error())
				hasErrors = true
			}
			img, err = resizeImage(w, h, img)
			if err != nil {
				log.Errorf("unable to resize image %s, err=%s", imagePath, err.Error())
				hasErrors = true
			}
			if addText {
				log.Infof("Adding watermark %s to image %s", wm[0].Text, imagePath)
				img, err = addWatermark(img, wm[0])
				if err != nil {
					log.Errorf("unable to add watermark %s to image %s. err=%v", wm[0].Text, imagePath, err.Error())
				}
			}
			imageName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			err = saveImage(dstPath, fmt.Sprintf("%s_%dx%d.jpg", imageName, w, h), img)
			if err != nil {
				log.Errorf("unable to save image %s. err=%s", path.Join(dstPath, fmt.Sprintf("%s_%dx%d.jpg", imageName, w, h)), err.Error())
				hasErrors = true
			}
		}
		return nil
	})
	if hasErrors {
		return fmt.Errorf("heck log for errors")
	}

	return nil
}

func resizeImage(w, h uint, img image.Image) (image.Image, error) {
	m := resize.Resize(w, h, img, resize.Lanczos3)
	return m, nil
}

func saveImage(filePath, fileName string, img image.Image) error {
	log.Infof("saving image %s", path.Join(filePath, fileName))
	return gg.SaveJPG(path.Join(filePath, fileName), img, 90)
}

func addWatermark(img image.Image, wm WaterMark) (image.Image, error) {
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(img, 0, 0)

	//TODO: Fixme this is dirty and only work for macos
	fontPath, err := findfont.Find("arial.ttf")
	if err != nil {
		fontPath, err = findfont.Find("Arial.ttf")
		if err != nil {
			return nil, fmt.Errorf("unable to load font. err=%s", err.Error())
		}
	}
	if err := dc.LoadFontFace(fontPath, wm.Size); err != nil {
		return nil, fmt.Errorf("unable to load font. err=%s", err.Error())
	}

	// Position on the image
	x := float64(imgWidth) - 200
	y := float64(imgHeight) - 20
	maxWidth := float64(imgWidth) - 60.0

	dc.SetColor(wm.Color)
	dc.DrawStringWrapped(wm.Text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	return dc.Image(), nil
}
