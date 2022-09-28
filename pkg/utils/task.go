package utils

import "image/color"

type Task struct {
	Path     string
	Name     string
	SavePath string
	Resize   struct {
		Enabled       bool
		Width, Height uint
	}
	Watermark struct {
		Enabled bool
		Color   color.Gray16
		Size    float64
		Text    string
	}
}
