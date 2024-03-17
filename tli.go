package main

import (
	"flag"
	"fmt"
	"github.com/DebuggerAndrzej/tli/ui"
)

func main() {
	filePath := flag.String("path", "", "Full path to log file")
	logFormat := flag.String("format", "", "log format")
	flag.Parse()

	if *filePath != "" {
		ui.InitTui(*filePath, *logFormat)
	} else {
		fmt.Println("Next time provide a file to open")
	}
}
