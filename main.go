package main

import (
	"context"
	"github.com/vfoucault/goPhoto/config"
	"github.com/vfoucault/goPhoto/logger"
	"github.com/vfoucault/goPhoto/photo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfg = new(config.Config)
)

type Photo struct {
	Name      string
	Path      string
	DateTaken time.Time
	DestPath  string
}

//
//func getUniqDestDirectories(colPhotos []Photo) map[string]int {
//	return_map := make(map[string]int)
//	for _, photo := range colPhotos {
//		_, ok := return_map[photo.DestPath]
//		if ok {
//			return_map[photo.DestPath] = return_map[photo.DestPath] + 1
//		} else {
//			return_map[photo.DestPath] = 1
//		}
//	}
//	return return_map
//}
//
//func copyImage(photo Photo) {
//	//srcFileInfo := os.FileInfo(photo.Path + "/" + photo.Name)
//	reader, _ := os.Open(photo.Path + "/" + photo.Name)
//
//	writer, _ := os.Create(photo.DestPath + "/" + photo.Name)
//	defer reader.Close()
//	defer writer.Close()
//	io.Copy(writer, reader)
//	writer.Sync()
//
//}

func main() {
	start := time.Now()
	cfg.Init()
	cfg.PrintConfig()

	ctx := context.Background()
	copier := photo.NewCopier(cfg, ctx)

	// Handle stop and more
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go handleSignals(signalChannel, copier)

	logger.Debugf("Will have to copy %v pictures", len(copier.Photos))

	if err := copier.Start(); err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	copier.Wait()

	elapsed := time.Since(start)
	logger.Infof("End of program. Took %v", elapsed)
}

//func handleSignals(signChannel chan os.Signal, processorManager *models.ProcessorManager) {
func handleSignals(signChannel chan os.Signal, copier *photo.Copier) {
	logger.Debugf("Running signal handler")
	// Handle stop and more
	for {
		sig := <-signChannel
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			logger.Infof("Received signal %v", sig)
			logger.Infof("Shutting down application. Kill running/pending jobs")
			copier.Stop()
			return
		default:
			logger.Errorf("Unable to handle signal %v", sig.String())
		}
	}

}
