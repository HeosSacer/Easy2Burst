package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"log"
)

// handleMessages handles messages
func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "exitEvent":
		commandCh <- "stopWallet"
		w.Close()
	case "startEvent":
		commandCh <- "startWallet"
	case "stopEvent":
		commandCh <- "stopWallet"
	case "bootstrapEvent":
		commandCh <- "bootstrapChain"
	default:
		log.Printf("Received unhandled message from ui %s", m.Name)
	}
	return
}
