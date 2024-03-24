package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DebuggerAndrzej/tli/ui"
)

func main() {
	logFormat := flag.String("f", "M", "log format")
	flag.Parse()
	filePath := flag.Arg(0)

	pipedInput := getPipedInput()

	if filePath == "" && pipedInput == "" {
		exitWithMessage("Please provide a way to get logs!")
	}
	if filePath != "" && pipedInput != "" {
		exitWithMessage("Can't handle both ways of passing inputs. At least for now...")
	}

	ui.InitTui(filePath, *logFormat, pipedInput)
}

func getPipedInput() string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return ""
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			exitWithMessage(fmt.Sprintf("Error getting input: %s", err))
		}
	}

	return b.String()
}

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}
