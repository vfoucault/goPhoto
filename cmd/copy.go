package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vfoucault/goPhoto/pkg/config"
	"github.com/vfoucault/goPhoto/pkg/photo"
)

var (
	srcDirectory   string
	dstDirectory   string
	dstFileFormat  string
	copyNoRecurse  bool
	copyNumWorkers int
)

// cmdAwsDelete delete ACM certificates
var cmdCopyPhoto = &cobra.Command{
	Use:     "copy",
	Short:   "Copy photos",
	Long:    "Copy photos from source to destination",
	Example: ``,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{
			DestFileFormat:  dstFileFormat,
			DestDirectory:   dstDirectory,
			SourceDirectory: srcDirectory,
			NoRecurse:       copyNoRecurse,
			Workers:         copyNumWorkers,
		}
		photo.RunCopier(cfg)
	},
}

func copyInit() {

	cmdCopyPhoto.PersistentFlags().StringVarP(&srcDirectory, "src", "s", ".", "Source directory")
	cmdCopyPhoto.PersistentFlags().StringVarP(&dstDirectory, "dst", "d", ".", "Destination directory")
	cmdCopyPhoto.MarkPersistentFlagRequired("src")
	cmdCopyPhoto.MarkPersistentFlagRequired("dst")
	cmdCopyPhoto.PersistentFlags().StringVarP(&dstFileFormat, "format", "", "2006/2006-01-02", "Destination directory format")
	cmdCopyPhoto.PersistentFlags().BoolVarP(&copyNoRecurse, "no-recurse", "", false, "Don't search recursively for photos")
	cmdCopyPhoto.PersistentFlags().IntVarP(&copyNumWorkers, "num-workers", "", runtime.NumCPU(), "number of workers. Default to runtime.NumCPU()")

	rootCmd.AddCommand(cmdCopyPhoto)

}
