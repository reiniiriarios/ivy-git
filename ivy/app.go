package ivy

import (
	"context"
	"ivy-git/ivy/git"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                  context.Context
	RepoSaveData         RepoSaveData
	Settings             Settings
	AppData              AppData
	WailsConfig          WailsConfig
	Git                  git.Git
	CurrentHash          string
	ShowRefAll           string
	UncommittedDiff      string
	UntrackedFiles       string
	RemoteDiff           string
	StagedDiff           string
	WatcherSemiSemaphore uint64
}

type WailsConfig struct {
	Name   string            `json:"name"`
	Slug   string            `json:"slug"`
	Bundle string            `json:"bundle"`
	Author wailsConfigAuthor `json:"author"`
	Info   wailsConfigInfo   `json:"info"`
}

type wailsConfigAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type wailsConfigInfo struct {
	Version string `json:"productVersion"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	// Context for app.
	a.ctx = ctx

	// Debug stats
	_, exists := os.LookupEnv("IVY_GIT_DEBUG")
	if exists {
		go a.statLoop()
	}

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
	if a.Git.Repo.Main == "" || !a.Git.BranchExists(a.Git.Repo.Main) {
		a.Git.Repo.Main = a.Git.NameOfMainBranch()
		r := a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo]
		r.Main = a.Git.Repo.Main
		a.RepoSaveData.Repos[a.RepoSaveData.CurrentRepo] = r
		a.saveRepoData()
	}
}

// called when the DOM is ready
func (a *App) Domready(ctx context.Context) {
	go a.watcher()
}
