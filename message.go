package main

import (
	"os"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"log"
)

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "exitEvent":
		w.Close()
		os.Exit(0)
	case "startEvent":
		//TODO
	case "stopEvent":
		//TODO
	case "bootstrapEvent":
		//TODO
	default:
		log.Printf("Received unhandled message from ui %s", m.Name)
	}
	return
}
