package main

import (
	"bytes"
	"log"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main() {

	startUI()
}