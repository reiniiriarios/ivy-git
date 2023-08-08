package main

import (
	"embed"
	"encoding/json"
	"fmt"

	"ivy-git/ivy"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

//go:embed wails.json
var wailsJson []byte

func main() {
	// Create an instance of the app structure
	app := ivy.NewApp()

	// Add config from wails.json for context
	app.WailsConfig = getWailsConfig()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Ivy Git",
		Width:     1280,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 600,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: appBgColor(),
		OnStartup:        app.Startup,
		OnDomReady:       app.Domready,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
		},
		Windows: &windows.Options{
			WindowIsTranslucent:  true,
			WebviewIsTransparent: true,
		},
		Linux: &linux.Options{
			Icon:                appIcon,
			WindowIsTranslucent: false,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
	})

	if err != nil {
		println("Error (main)", err.Error())
	}
}

func getWailsConfig() ivy.WailsConfig {
	var cfg ivy.WailsConfig
	if err := json.Unmarshal(wailsJson, &cfg); err != nil {
		fmt.Println(err.Error())
	}

	return cfg
}
