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
	ShowRefAll           string
	UncommittedDiff      string
	UntrackedFiles       string
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
	// Context for app.
	a.ctx = ctx

	// Load yaml configs.
	a.loadConfig()

	// Set app data.
	if a.AppData.WindowWidth > 1024 && a.AppData.WindowHeight > 600 {
		runtime.WindowSetSize(a.ctx, a.AppData.WindowWidth, a.AppData.WindowHeight)
	}

	// Set git data.
	a.Git = git.Git{
		AppCtx: a.ctx,
		Repo:   a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo],
	}
	// If main branch not found, check again.
	if !a.Git.BranchExists(a.Git.Repo.Main) {
		a.Git.Repo.Main = a.Git.NameOfMainBranch()
	}
}

// called when the DOM is ready
func (a *App) domready(ctx context.Context) {
	go a.watcher()
}
