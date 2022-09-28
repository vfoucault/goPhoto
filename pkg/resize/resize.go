package resize

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

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	"github.com/vfoucault/goPhoto/pkg/utils"
	"github.com/vfoucault/goPhoto/pkg/watermark"
)

func PhotoResize(w, h uint, srcPath, dstPath string, addText bool, wm ...watermark.WaterMark) error {
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
				Resize: struct {
					Enabled       bool
					Width, Height uint
				}{
					Enabled: true,
					Width:   w,
					Height:  h,
				},
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
	log.Infof("resizing %s...", task.Name)
	imagePath := path.Join(task.Path, task.Name)
	img, err := gg.LoadImage(imagePath)
	var hasErrors bool
	if err != nil {
		log.Errorf("unable to load image %s. err=%w", imagePath, err.Error())
		hasErrors = true
	}
	img, err = resizeImage(task.Resize.Width, task.Resize.Height, img)
	if err != nil {
		log.Errorf("unable to resize image %s, err=%s", imagePath, err.Error())
		hasErrors = true
	}
	if task.Watermark.Enabled {
		log.Infof("Adding watermark %s to image %s", task.Watermark.Text, imagePath)
		img, err = watermark.AddWatermark(img, task.Watermark.Text, task.Watermark.Color, task.Watermark.Size)
		if err != nil {
			log.Errorf("unable to add watermark %s to image %s. err=%v", task.Watermark.Text, imagePath, err.Error())
		}
	}
	imageName := strings.TrimSuffix(task.Name, filepath.Ext(task.Name))
	err = utils.SaveImage(task.SavePath, fmt.Sprintf("%s_%dx%d.jpg", imageName, task.Resize.Width, task.Resize.Height), img)
	if err != nil {
		log.Errorf("unable to save image %s. err=%s", path.Join(task.SavePath, fmt.Sprintf("%s_%dx%d.jpg", imageName, task.Resize.Width, task.Resize.Height)), err.Error())
		hasErrors = true
	}
	return hasErrors
}

func resizeImage(w, h uint, img image.Image) (image.Image, error) {
	m := resize.Resize(w, h, img, resize.Lanczos3)
	return m, nil
}
