package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	RepoSaveData RepoSaveData
	Settings     Settings
	AppData      AppData
}

type GenericResponse struct {
	Response string
	Message  string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadConfig()
	if a.AppData.WindowWidth > 1024 && a.AppData.WindowHeight > 600 {
		runtime.WindowSetSize(a.ctx, a.AppData.WindowWidth, a.AppData.WindowHeight)
	}
}

func (a *App) domready(ctx context.Context) {
	// ...
}
