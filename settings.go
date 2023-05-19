package main

import (
	"errors"
	"io/ioutil"
	"ivy-git/git"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

const NAME = "Ivy Git"
const SHORT_NAME = "ivy-git"
const BUNDLE = "me.reinii.ivy-git"
const VERSION = "0.0.1"

type RepoSaveData struct {
	CurrentRepo string
	Repos       map[string]git.Repo
}

type Settings struct {
	Version                      string
	DisplayCommitSignatureInList bool
}

type AppData struct {
	WindowWidth  int
	WindowHeight int
}

// Load configuration yaml for app.
func (a *App) loadConfig() {
	rp := filepath.Join(a.settingsLocationLocal(), "repos.yaml")
	repo_data := a.initConfigFile(rp)
	var repos RepoSaveData
	err := yaml.Unmarshal(repo_data, &repos)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	a.RepoSaveData = repos

	sp := filepath.Join(a.settingsLocationRoaming(), "settings.yaml")
	settings_data := a.initConfigFile(sp)
	var settings Settings
	err = yaml.Unmarshal(settings_data, &settings)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	a.Settings = settings

	dp := filepath.Join(a.settingsLocationLocal(), "appdata.yaml")
	app_data := a.initConfigFile(dp)
	var data AppData
	err = yaml.Unmarshal(app_data, &data)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	a.AppData = data
}

// Save repo data to config file.
func (a *App) saveRepoData() {
	data, err := yaml.Marshal(&a.RepoSaveData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	rp := filepath.Join(a.settingsLocationLocal(), "repos.yaml")
	err2 := os.WriteFile(rp, []byte(data), 0644)
	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}
}

// Save settings to config file.
func (a *App) saveSettings() {
	a.Settings.Version = VERSION
	data, err := yaml.Marshal(&a.Settings)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	rp := filepath.Join(a.settingsLocationRoaming(), "settings.yaml")
	err2 := os.WriteFile(rp, []byte(data), 0644)
	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}
}

// Save app data to config file.
func (a *App) saveData() {
	data, err := yaml.Marshal(&a.AppData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	rp := filepath.Join(a.settingsLocationLocal(), "appdata.yaml")
	err2 := os.WriteFile(rp, []byte(data), 0644)
	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}
}

// Return file data, create empty file if not found.
func (a *App) initConfigFile(file string) []byte {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, e := os.Create(file)
		if e != nil {
			runtime.LogError(a.ctx, e.Error())
		}
		defer f.Close()
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	return contents
}

// Get location for local settings,
// specific to the machine the app is running on.
func (a *App) settingsLocationLocal() string {
	// Darwin
	//   ~/Library/Caches
	// Windows
	//   %LocalAppData%
	//   C:\Users\YourUser\AppData\Local
	// Unix
	//   $XDG_CACHE_HOME or $HOME/.cache
	d, err := os.UserCacheDir()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	dir := filepath.Join(d, a.settingsDirName())
	a.initDir(dir)

	return dir
}

// Get location for roaming settings,
// specific to the user.
func (a *App) settingsLocationRoaming() string {
	// Darwin
	//   ~/Library/Application Support/
	// Windows
	//   %AppData%
	//   C:\Users\YourUser\AppData\Roaming
	// Unix
	//   $XDG_CONFIG_HOME or $HOME/.config
	d, err := os.UserConfigDir()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	dir := filepath.Join(d, a.settingsDirName())
	a.initDir(dir)

	return dir
}

// Create a dir if it doesn't exist.
func (a *App) initDir(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
		}
	}
}

// Get name of directory to save settings under.
func (a *App) settingsDirName() string {
	if runtime.Environment(a.ctx).Platform == "windows" {
		// Windows convention is the name of the app.
		return NAME
	} else if runtime.Environment(a.ctx).Platform == "darwin" {
		// Darwin convention is the bundle id.
		return BUNDLE
	}
	// For unix, use the shortname.
	return SHORT_NAME
}

// When window resizes, save that info.
func (a *App) ResizeWindow() bool {
	w, h := runtime.WindowGetSize(a.ctx)
	a.AppData.WindowWidth = w
	a.AppData.WindowHeight = h
	a.saveData()
	return true
}

// Get settings in frontend.
func (a *App) GetSettings() Settings {
	return a.Settings
}
