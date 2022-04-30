package main

import (
	"github.com/vfoucault/goPhoto/cmd"
	"github.com/vfoucault/goPhoto/pkg/logger"
)

func main() {
	logger.SetupLog()
	cmd.Execute()
}
