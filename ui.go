package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/HeosSacer/Easy2Burst/internal"
	"github.com/pkg/errors"
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

func startUI (chan internal.Status) {
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

