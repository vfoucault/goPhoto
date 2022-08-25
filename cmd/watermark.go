package cmd

import (
	"image/color"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vfoucault/goPhoto/pkg/watermark"
)

var (
	watermarkSrcDirectory   string
	watermarkDstDirectory   string
	watermarkWatermarkText  string
	watermarkWatermarkColor string
	watermarkWatermarkSize  float64
)

var cmdWatermark = &cobra.Command{
	Use:     "watermark",
	Short:   "Add watermark to pictures",
	Long:    "Watermark photos from source to destination",
	Example: ``,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// Watermark
		wm := watermark.WaterMark{Size: watermarkWatermarkSize, Text: watermarkWatermarkText}
		switch watermarkWatermarkColor {
		case "white":
			wm.Color = color.White
		case "black":
			wm.Color = color.Black
		default:
			log.Errorf("unable to process color %s. only white and black", watermarkWatermarkColor)
		}
		log.Infof("calling add watermark with %s, %s, %s", watermarkSrcDirectory, watermarkDstDirectory, len(watermarkWatermarkText))
		err := watermark.AddWatermarkToImage(watermarkSrcDirectory, watermarkDstDirectory, len(watermarkWatermarkText) > 0, wm)
		if err != nil {
			log.Errorf("unable to start image resize processing. err=%v", err.Error())
		}
	},
}

func watermarkInit() {

	cmdWatermark.PersistentFlags().StringVarP(&watermarkSrcDirectory, "src", "s", ".", "Source directory")
	cmdWatermark.PersistentFlags().StringVarP(&watermarkDstDirectory, "dst", "d", ".", "Destination directory")
	cmdWatermark.MarkPersistentFlagRequired("src")
	cmdWatermark.MarkPersistentFlagRequired("dst")
	cmdWatermark.PersistentFlags().StringVarP(&watermarkWatermarkText, "watermark", "", "", "Watermark text")
	cmdWatermark.PersistentFlags().StringVarP(&watermarkWatermarkColor, "watermark-color", "", "white", "Watermark color")
	cmdWatermark.PersistentFlags().Float64VarP(&watermarkWatermarkSize, "watermark-size", "", 25, "Watermark color (black / white)")

	rootCmd.AddCommand(cmdWatermark)

}
