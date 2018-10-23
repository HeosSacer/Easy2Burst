package internal

import (
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"github.com/asticode/go-astilectron"
	"strings"
	"os"
	"log"
)

var(
	Log *log.Logger
)

func StartUiManager(statusCh chan Status, command chan string, window *astilectron.Window){
	os.Remove(toolPath + "ui.log")
	logFile, _ := os.OpenFile(toolPath + "ui.log", os.O_WRONLY|os.O_CREATE, 0644)
	Log = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	ControlLoop:
		for{
			status := <-statusCh
			switch status.Name {
			case "fatalError":
				break ControlLoop
			default:
				sendMessageToUi(window, status.Name, status.Message, status.Progress, status.Size)
			}
		}
	defer logFile.Close()
	window.Close()
}

func sendMessageToUi(w *astilectron.Window, msgName string, payload ...string){
	payloadMsg := ""
	if len(payload) > 0{
		payloadMsg = strings.Join(payload,";")
	}
	Log.Print("Sending " + msgName + " PL: " + payloadMsg)
	if err := bootstrap.SendMessage(w, msgName, payloadMsg); err != nil {
		astilog.Error(errors.Wrap(err, "sending " + msgName + " with payload " + payloadMsg + " failed"))
		Log.Fatal(err)
	}
}

