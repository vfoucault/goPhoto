package watermark

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/flopp/go-findfont"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"
	"github.com/vfoucault/goPhoto/pkg/utils"
)

type WaterMark struct {
	Color color.Gray16
	Size  float64
	Text  string
}

func AddWatermarkToImage(srcPath, dstPath string, addText bool, wm ...WaterMark) error {
	// list all images
	tasks := make(chan *utils.Task, runtime.NumCPU())
	// Launch workers
	wg := &sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				task := <-tasks
				//log.Infof("Will have to run task for task %s", task.Name)
				hasErrors := ProcessTask(task)
				if hasErrors {
					log.Errorf("processing %s led to errors. check logs.", task.Name)
				}
				wg.Done()
			}
		}()
	}
	filepath.Walk(srcPath, func(aPath string, f os.FileInfo, _ error) error {
		if utils.IsImage(f) {
			// create Task struct
			task := &utils.Task{
				Path:     strings.TrimSuffix(aPath, f.Name()),
				Name:     f.Name(),
				SavePath: dstPath,
				Watermark: struct {
					Enabled bool
					Color   color.Gray16
					Size    float64
					Text    string
				}{
					Enabled: addText,
					Color:   wm[0].Color,
					Size:    wm[0].Size,
					Text:    wm[0].Text,
				},
			}
			tasks <- task
			wg.Add(1)
		}
		return nil
	})
	wg.Wait()
	return nil
}

func ProcessTask(task *utils.Task) bool {
	imagePath := path.Join(task.Path, task.Name)
	img, err := gg.LoadImage(imagePath)
	var hasErrors bool
	if err != nil {
		log.Errorf("unable to load image %s. err=%w", imagePath, err.Error())
		hasErrors = true
	}
	if task.Watermark.Enabled {
		log.Infof("Adding watermark %s to image %s", task.Watermark.Text, imagePath)
		img, err = AddWatermark(img, task.Watermark.Text, task.Watermark.Color, task.Watermark.Size)
		if err != nil {
			log.Errorf("unable to add watermark %s to image %s. err=%v", task.Watermark.Text, imagePath, err.Error())
		}
	}
	imageName := strings.TrimSuffix(task.Name, filepath.Ext(task.Name))
	err = utils.SaveImage(task.SavePath, fmt.Sprintf("%s.jpg", imageName), img)
	if err != nil {
		log.Errorf("unable to save image %s. err=%s", path.Join(task.SavePath, fmt.Sprintf("%s_%dx%d.jpg", imageName, task.Resize.Width, task.Resize.Height)), err.Error())
		hasErrors = true
	}
	return hasErrors
}

func AddWatermark(img image.Image, text string, c color.Gray16, size float64) (image.Image, error) {
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(img, 0, 0)

	fontPath, err := findfont.Find("arial.ttf")
	if err != nil {
		fontPath, err = findfont.Find("Arial.ttf")
		if err != nil {
			return nil, fmt.Errorf("unable to load font. err=%s", err.Error())
		}
	}
	if err := dc.LoadFontFace(fontPath, size); err != nil {
		return nil, fmt.Errorf("unable to load font. err=%s", err.Error())
	}

	// Position on the image
	x := float64(imgWidth) - 200
	y := float64(imgHeight) - 20
	maxWidth := float64(imgWidth) - 60.0

	dc.SetColor(c)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	return dc.Image(), nil
}
