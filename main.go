package main

import (
	"bytes"
	"log"
	"github.com/HeosSacer/Easy2Burst/internal"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	statusCh = make(chan internal.Status)
)

func main() {
	go internal.CheckTools(statusCh)
	startUI(statusCh)
}