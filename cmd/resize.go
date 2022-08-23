package cmd

import (
	"image/color"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vfoucault/goPhoto/pkg/resize"
)

var (
	resizeSrcDirectory   string
	resizeDstDirectory   string
	resizeWatermarkText  string
	resizeWatermarkColor string
	resizeWatermarkSize  float64
	resizeSize           string
)

// cmdAwsDelete delete ACM certificates
var cmdResize = &cobra.Command{
	Use:     "resize",
	Short:   "resize photos",
	Long:    "Resize photos from source to destination",
	Example: ``,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// define target size it not
		sizeList := strings.Split(resizeSize, "x")
		if len(sizeList) != 2 {
			log.Errorf("unable to parse %s as size if format XxY", resizeSize)
			return
		}
		w, err := strconv.ParseUint(sizeList[0], 10, 32)
		if err != nil {
			log.Errorf("unable to parse width %s. err=%v", sizeList[0], err.Error())
		}
		h, err := strconv.ParseUint(sizeList[1], 10, 32)
		if err != nil {
			log.Errorf("unable to parse width %s. err=%v", sizeList[0], err.Error())
		}

		// Watermark
		wm := resize.WaterMark{Size: resizeWatermarkSize, Text: resizeWatermarkText}
		switch resizeWatermarkColor {
		case "white":
			wm.Color = color.White
		case "black":
			wm.Color = color.Black
		default:
			log.Errorf("unable to process color %s. only white and black", resizeWatermarkColor)
		}

		err = resize.PhotoResize(uint(w), uint(h), resizeSrcDirectory, resizeDstDirectory, len(resizeWatermarkText) > 0, wm)
		if err != nil {
			log.Errorf("unable to start image resize processing. err=%v", err.Error())
		}
	},
}

func resizeInit() {

	cmdResize.PersistentFlags().StringVarP(&resizeSrcDirectory, "src", "s", ".", "Source directory")
	cmdResize.PersistentFlags().StringVarP(&resizeDstDirectory, "dst", "d", ".", "Destination directory")
	cmdResize.MarkPersistentFlagRequired("src")
	cmdResize.MarkPersistentFlagRequired("dst")
	cmdResize.PersistentFlags().StringVarP(&resizeSize, "size", "r", "1600x1064", "target size")
	cmdResize.PersistentFlags().StringVarP(&resizeWatermarkText, "watermark", "", "", "Watermark text")
	cmdResize.PersistentFlags().StringVarP(&resizeWatermarkColor, "watermark-color", "", "white", "Watermark color")
	cmdResize.PersistentFlags().Float64VarP(&resizeWatermarkSize, "watermark-size", "", 25, "Watermark color (black / white)")

	rootCmd.AddCommand(cmdResize)

}
