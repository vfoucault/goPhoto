package main

import (
	"io/ioutil"
	"log"
	"flag"
	"os"
	"github.com/rwcarlsen/goexif/exif"
	"time"
	"io"
	"path/filepath"
	"regexp"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	//"go/constant"
)

const fileMode = 0755

var (
	Trace   *log.Logger
	Info    *log.Logger
	Debug   *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	regexImage *regexp.Regexp
	destFormat string
	destDirectory string
)

func Init(
traceHandle io.Writer,
infoHandle io.Writer,
warningHandle io.Writer,
errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Debug = log.New(infoHandle,
		"DEBUG: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	regexImage, _ = regexp.Compile(`^.+(?i)(jpe?g|gif|png|tiff)$`)

}

type Photo struct {
	Name      string
	Path      string
	DateTaken time.Time
	DestPath  string
}

//func setDestPath(photo Photo, destDirectory string) (Photo) {

//}

func GetDateTaken(photo Photo) (Photo) {
	f, err := os.Open(photo.Path + "/" + photo.Name)
	if err != nil {
		log.Fatal(err)
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	photo.DateTaken, _ = x.DateTime()
	photo.DestPath = destDirectory + "/" + photo.DateTaken.Format(destFormat)
	return photo

}

func isImage(file os.FileInfo) (bool) {
	if !file.IsDir() {
		if regexImage.MatchString(file.Name()) {
			return true

		} else {
			return false
		}
	}
	return false
}

func getImages(directory string, noRecursive bool) ([]Photo) {

	var colPhoto []Photo
	var files []os.FileInfo
	switch noRecursive  {
	case false:
		filepath.Walk(directory, func(path string, f os.FileInfo, _ error) error {
			if isImage(f) {
				photo := GetDateTaken(Photo{Name:f.Name(), Path:filepath.Dir(path)})
				colPhoto = append(colPhoto, photo)
			}
			return nil
		})
	default:
		files, _ = ioutil.ReadDir(directory)
		for _, f := range files {
			if isImage(f) {
				photo := GetDateTaken(Photo{Name:f.Name(), Path:directory})
				colPhoto = append(colPhoto, photo)
			}
		}
	}
	return colPhoto
}

func createDirectories(directory string) {
	err := os.MkdirAll(directory, fileMode)
	if err != nil {
		log.Fatal("Something went wront creating ", directory, err)
	}

}

func getUniqDestDirectories(colPhotos []Photo) (map[string]int) {
	return_map := make(map[string]int)
	for _, photo := range colPhotos {
		_, ok := return_map[photo.DestPath]
		if ok {
			return_map[photo.DestPath] = return_map[photo.DestPath] + 1
		} else {
			return_map[photo.DestPath] = 1
		}
	}
	return return_map
}

func copyImage(photo Photo) {
	//srcFileInfo := os.FileInfo(photo.Path + "/" + photo.Name)
	reader, _ := os.Open(photo.Path + "/" + photo.Name)

	writer, _ := os.Create(photo.DestPath + "/" + photo.Name)
	defer reader.Close()
	defer writer.Close()
	io.Copy(writer, reader)

}

func main() {

	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	fmt.Println(`

       .__            __         .___                              __
______ |  |__   _____/  |_  ____ |   | _____ ______   ____________/  |_
\____ \|  |  \ /  _ \   __\/  _ \|   |/     \\____ \ /  _ \_  __ \   __\
|  |_> >   Y  (  <_> )  | (  <_> )   |  Y Y  \  |_> >  <_> )  | \/|  |
|   __/|___|  /\____/|__|  \____/|___|__|_|  /   __/ \____/|__|   |__|
|__|        \/                             \/|__|

	`)

	srcFlag := flag.String("source", "./", "the source directory")
	noRecurseFlag := flag.Bool("no-recurse", false, "Don't search recursively ?")
	dstFlag := flag.String("destination", "", "the destination directory")
	dstFormatFlag := flag.String("format", "2006/2006-01-02", "the destination format")
	//moveFlag := flag.Bool("move", false, "move files ?")

	flag.Parse()

	destDirectory = *dstFlag
	destFormat = *dstFormatFlag

	Info.Println("Running with source =", *srcFlag)
	Info.Println("Running with destination =", *dstFlag)
	Info.Println("Running with destination format =", *dstFormatFlag)

	colPhotos := getImages(*srcFlag, *noRecurseFlag)

	Info.Println("Number of images files :", len(colPhotos))

	mapDestFolders := getUniqDestDirectories(colPhotos)

	for directory := range mapDestFolders {
		if _, err := os.Stat(directory); err != nil {
			if os.IsNotExist(err) {
				Info.Println("Creating directory :", directory)
				createDirectories(directory)
			}
		}
	}

	for _, photo := range colPhotos {
		Info.Println("Copying ", photo.Name, "to", photo.DestPath)

		copyImage(photo)
	}

	Info.Println("End of program. Took", )

}