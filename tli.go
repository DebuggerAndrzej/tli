package main

import (
	"flag"
	"fmt"
	"tli/ui"
)

func main() {
	filePath := flag.String("path", "", "Full path to config file")
	flag.Parse()

	if *filePath != "" {
		ui.InitTui(*filePath)
	} else {
		fmt.Println("Next time provide a file to open")
	}
}
