package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "BulletinBoard",
		Width:             100,
		Height:            60,
		Assets:            assets,
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 0},
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       true,
		AlwaysOnTop:       true,
		HideWindowOnClose: true,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		CSSDragProperty:   "--wails-draggable",
		CSSDragValue:      "drag",
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "BulletinBoard",
				Message: "© 2022 Richard Guay <raguay@customct.com>",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
