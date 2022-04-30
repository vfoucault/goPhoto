package photo

import (
	"context"
	"sync"
	"testing"

	"github.com/schollz/progressbar/v3"
	"github.com/vfoucault/goPhoto/pkg/config"
)

func TestCopier_IncrementStats(t *testing.T) {
	type fields struct {
		Config      *config.Config
		Photos      []*Photo
		Stats       int
		StatsMutex  sync.Mutex
		Workers     []*Worker
		CopyQueue   chan *Photo
		Context     context.Context
		CancelFunc  context.CancelFunc
		Wg          sync.WaitGroup
		ProgressBar *progressbar.ProgressBar
	}
	tests := []struct {
		name   string
		want   int
		fields fields
	}{
		{
			name: "Should increment stats",
			fields: fields{
				Stats:      0,
				StatsMutex: sync.Mutex{},
			},
			want: 1,
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
			c.IncrementStats()
			if c.Stats != tt.want {
				t.Errorf("IncrementStats(): got %d want %d", c.Stats, tt.want)
			}
		})
	}
}
