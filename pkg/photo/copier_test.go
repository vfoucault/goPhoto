package photo

import (
	"context"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/vfoucault/goPhoto/pkg/config"
)

func TestCopier_incrementStats(t *testing.T) {
	type fields struct {
		Config *config.Config
		Photos []*Photo
		Stats  struct {
			Count int
			Size  int64
		}
		StatsMutex  sync.Mutex
		Workers     []*Worker
		CopyQueue   chan *Photo
		Context     context.Context
		CancelFunc  context.CancelFunc
		Wg          sync.WaitGroup
		ProgressBar *progressbar.ProgressBar
	}
	type args struct {
		size int64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantSize  int64
	}{
		{
			name: "Should increment the stats",
			fields: fields{
				Stats: struct {
					Count int
					Size  int64
				}{},
				StatsMutex: sync.Mutex{},
			},
			args: args{
				size: 128,
			},
			wantCount: 1,
			wantSize:  128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Copier{
				Config:      tt.fields.Config,
				Photos:      tt.fields.Photos,
				Stats:       tt.fields.Stats,
				StatsMutex:  tt.fields.StatsMutex,
				Workers:     tt.fields.Workers,
				CopyQueue:   tt.fields.CopyQueue,
				Context:     tt.fields.Context,
				CancelFunc:  tt.fields.CancelFunc,
				Wg:          tt.fields.Wg,
				ProgressBar: tt.fields.ProgressBar,
			}
			c.incrementStats(tt.args.size)
			if c.Stats.Size != tt.wantSize {
				t.Errorf("incrementStats() got size=%d want size=%d", c.Stats.Size, tt.wantSize)
			}
			if c.Stats.Count != tt.wantCount {
				t.Errorf("incrementStats() got count=%d want count=%d", c.Stats.Count, tt.wantCount)
			}
		})
	}
}

func TestCopier_CreateDestDirs(t *testing.T) {
	tmpDestDir, err := os.MkdirTemp(os.TempDir(), "goPhotos_tests")
	if err != nil {
		t.Errorf("unable to create temp directory. err=%v", err.Error())
	}
	defer os.RemoveAll(tmpDestDir)

	type fields struct {
		Photos []*Photo
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Should Create all needed directories",
			fields: fields{
				Photos: []*Photo{{
					Path:     path.Join(tmpDestDir, "some_path"),
					FileName: "img001.jpg",
					Copier: &Copier{
						Config: &config.Config{
							DestDirectory:  tmpDestDir,
							DestFileFormat: "2006-01-02",
						},
					},
					//2022-04-30
					DateTaken: time.Unix(1651276800, 0),
					Atime:     time.Time{},
					Ctime:     time.Time{},
					Mtime:     time.Time{},
				}, {
					Path:     path.Join(tmpDestDir, "some_path"),
					FileName: "img002.jpg",
					Copier: &Copier{
						Config: &config.Config{
							DestDirectory:  tmpDestDir,
							DestFileFormat: "2006-01-02",
						},
					},
					//2022-03-29
					DateTaken: time.Unix(1648512000, 0),
					Atime:     time.Time{},
					Ctime:     time.Time{},
					Mtime:     time.Time{},
				}},
			},
			want: []string{"2022-04-30", "2022-03-29"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Copier{
				Photos: tt.fields.Photos,
			}
			c.CreateDestDirs()
			// check that needed directories where created
			for _, dir := range tt.want {
				if _, err := os.Stat(path.Join(tmpDestDir, dir)); err != nil {
					t.Errorf("CreateDestDirs(): directory %s not found", path.Join(tmpDestDir, dir))
				}
			}
		})
	}
}

//func TestCopier_Search(t *testing.T) {
//	tmpDestDir, err := os.MkdirTemp(os.TempDir(), "goPhotos_tests")
//
//	if err != nil {
//		t.Errorf("unable to create temp directory. err=%v", err.Error())
//	}
//	//defer os.RemoveAll(tmpDestDir)
//	tmpfile1, _ := ioutil.TempFile(tmpDestDir, "testfile1.*.jpg")
//	tmpfile2, _ := ioutil.TempFile(tmpDestDir, "testfile2.*.jpg")
//
//	exifData1, _ := exiftool.NewExiftool()
//	defer exifData1.Close()
//	exifData1.
//	type fields struct {
//		Config *config.Config
//		Photos []*Photo
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "Should populate photo on non recursive",
//			fields: fields{
//				Config: &config.Config{
//					SourceDirectory: "",
//					NoRecurse:       false,
//					Verbose:         false,
//					Workers:         0,
//				},
//				Photos: nil,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Copier{
//				Config: tt.fields.Config,
//				Photos: tt.fields.Photos,
//			}
//			c.Search()
//		})
//	}
//}
