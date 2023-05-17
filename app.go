package main

import (
	"context"
	"ivy-git/git"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                  context.Context
	RepoSaveData         RepoSaveData
	Settings             Settings
	AppData              AppData
	Git                  git.Git
	CurrentHash          string
	UncommittedDiff      string
	RemoteDiff           string
	StagedDiff           string
	WatcherSemiSemaphore uint64
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
	a.Git = git.Git{
		Repo: a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo],
	}
}

// called when the DOM is ready
func (a *App) domready(ctx context.Context) {
	go a.watcher()
}
