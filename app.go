package main

import (
	"context"
	"runtime"
)

// App struct
type App struct {
	ctx          context.Context
	RepoSaveData RepoSaveData
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
	a.LoadYaml()
}

func (a *App) domready(ctx context.Context) {
	// ...
}

func (a *App) GoOs() string {
	return runtime.GOOS
}
