package main

import (
	"embed"
	"encoding/json"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/config.json
var configFile []byte

type Config struct {
	AppSideLength int `json:"APP_SIDE_LENGTH"`
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Parse config
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		println("Error parsing config:", err.Error())
		// Fallback defaults
		config.AppSideLength = 1024
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "goldenfox",
		Width:         config.AppSideLength,
		Height:        config.AppSideLength,
		DisableResize: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
