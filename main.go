package main

import (
	"bytes"
	"github.com/HeosSacer/Easy2Burst/internal"
	"log"
	"os"
)

var (
	buf       bytes.Buffer
	logger    = log.New(&buf, "logger: ", log.Lshortfile)
	statusCh  = make(chan internal.Status)
	commandCh = make(chan string)
)

func main() {
	//Check for args
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--beta":
			//TODO
		case "--update":
			pathOfOldBinary := os.Args[2]
			internal.UpdateBinary(pathOfOldBinary)
			os.Exit(1)
		}
	}
	go internal.CheckTools(statusCh, commandCh)
	startUI(statusCh, commandCh)
}
