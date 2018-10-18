package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/HeosSacer/Easy2Burst/internal"
	"github.com/pkg/errors"
	"time"
)

// Constants
const htmlAbout = `Easy2Burst Project.`

// Vars
var (
	AppName string
	BuiltAt string
	debug   = true
	w       *astilectron.Window
)

func startUI (statusCh chan internal.Status, commandCh chan string) {
	// Init
	flag.Parse()
	astilog.FlagInit()
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: debug,
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go internal.StartUiManager(statusCh, commandCh, w)
			if err := bootstrap.SendMessage(w, "starting", "0%;"); err != nil {
				astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
			}
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Frame: 			 astilectron.PtrBool(false),
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(829),
				Width:           astilectron.PtrInt(1259),
				TitleBarStyle: 	 astilectron.PtrStr("customButtonsOnHover"),
				WebPreferences: &astilectron.WebPreferences{
					DevTools: astilectron.PtrBool(true),
				},
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

