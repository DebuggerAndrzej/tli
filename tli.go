package main

import (
	"flag"
	"fmt"

	"github.com/DebuggerAndrzej/tli/ui"
)

func main() {
	logFormat := flag.String("f", "M", "log format")
	flag.Parse()
	filePath := flag.Arg(0)

	if filePath != "" {
		ui.InitTui(filePath, *logFormat)
	} else {
		fmt.Println("Next time provide a file to open")
	}
}
