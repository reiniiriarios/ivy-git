package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ivy-git/git"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

type RepoSaveData struct {
	CurrentRepo string
	Repos       map[string]git.Repo
}

type Settings struct {
	Version                      string
	Workflow                     string
	Theme                        string
	DateFormat                   uint8
	HighlightMainBranch          bool
	HighlightConventionalCommits bool
	DisplayCommitSignatureInList bool
	DisplayAvatars               bool
	BackgroundOpacity            uint8
	AutoFetch                    bool
}

type AppData struct {
	WindowWidth   int
	WindowHeight  int
	RecentRepoDir string
}

// Get default settings. These are immediately overwritten by the saved data in the yaml file.
func (a *App) getDefaultSettings() Settings {
	return Settings{
		Workflow:                     "merge",
		DateFormat:                   0,
		HighlightMainBranch:          true,
		HighlightConventionalCommits: false,
		DisplayCommitSignatureInList: false, // wip
		DisplayAvatars:               true,
		BackgroundOpacity:            DEFAULT_BG_OPACITY,
		AutoFetch:                    true,
	}
}

// Load configuration yaml for app.
func (a *App) loadConfig() {
	// Get repo data.
	rp := filepath.Join(a.settingsLocationLocal(), "repos.yaml")
	repo_data := a.initConfigFile(rp)
	var repos RepoSaveData
	err := yaml.Unmarshal(repo_data, &repos)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	// Validate repo data.
	if _, exists := repos.Repos[repos.CurrentRepo]; !exists {
		// If the current repo isn't in the list of repos (unlikely, but possible due to bugs).
		repos.CurrentRepo = ""
	} else if !a.Git.IsGitRepo(repos.Repos[repos.CurrentRepo].Directory) {
		// If the current repo isn't found or is no longer a git repo.
		repos.CurrentRepo = ""
	}
	if repos.Repos == nil {
		repos.Repos = make(map[string]git.Repo)
	}
	a.RepoSaveData = repos

	// Get settings.
	sp := filepath.Join(a.settingsLocationRoaming(), "settings.yaml")
	settings_data := a.initConfigFile(sp)
	settings := a.getDefaultSettings()
	err = yaml.Unmarshal(settings_data, &settings)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		settings = a.getDefaultSettings()
	}
	// Validate settings.
	if settings.BackgroundOpacity < 10 || settings.BackgroundOpacity > 100 {
		settings.BackgroundOpacity = DEFAULT_BG_OPACITY
	}
	// Call the save method to validate and correct outdated settings.
	a.saveSettings(settings)

	// Get app data.
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
func (a *App) saveSettings(new_settings Settings) error {
	// Always save the current version in settings.
	new_settings.Version = a.getWailsConfig().Info.Version

	// Only three viable options for this.
	if new_settings.Workflow != "squash" && new_settings.Workflow != "rebase" {
		new_settings.Workflow = "merge"
	}

	a.Settings = new_settings

	data, err := yaml.Marshal(&a.Settings)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	rp := filepath.Join(a.settingsLocationRoaming(), "settings.yaml")
	err = os.WriteFile(rp, []byte(data), 0644)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	return nil
}

// Handler for frontend.
func (a *App) SaveSettingsGui(new_settings Settings) DataResponse {
	err := a.saveSettings(new_settings)
	return dataResponse(err, false)
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

// After app data is saved, sometimes we emit to the frontend.
func (a *App) emitData() {
	runtime.EventsEmit(a.ctx, "appdata", a.AppData)
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
		return a.getWailsConfig().Name
	} else if runtime.Environment(a.ctx).Platform == "darwin" {
		// Darwin convention is the bundle id.
		return a.getWailsConfig().Bundle
	}
	// For unix, use the slug.
	return a.getWailsConfig().Slug
}

// Get settings in frontend.
func (a *App) GetSettings() Settings {
	return a.Settings
}

// Get app data in frontend.
func (a *App) GetAppData() AppData {
	return a.AppData
}

// When window resizes, save that info.
func (a *App) ResizeWindow() bool {
	w, h := runtime.WindowGetSize(a.ctx)
	a.AppData.WindowWidth = w
	a.AppData.WindowHeight = h
	a.saveData()
	return true
}

// Save most recent directory location for repos.
func (a *App) saveRecentRepoDir(dir string) bool {
	a.AppData.RecentRepoDir = dir
	a.saveData()
	a.emitData()
	return true
}

type DateFormat struct {
	Id      uint8
	Display string
	format  string
}

// Date format options.
// Do not change IDs.
func getDateFormatData() map[uint8]DateFormat {
	formats := make(map[uint8]DateFormat)
	formats[0] = DateFormat{
		Display: "12hr",
		format:  "Jan 2, 2006, 03:04:05 pm",
	}
	formats[1] = DateFormat{
		Display: "24hr",
		format:  "Jan 2, 2006, 15:04:05 pm",
	}
	formats[2] = DateFormat{
		Display: "12hr",
		format:  "2006-01-02, 03:04:05 pm",
	}
	formats[3] = DateFormat{
		Display: "24hr",
		format:  "2006-01-02, 15:04:05 pm",
	}
	formats[4] = DateFormat{
		Display: "12hr",
		format:  "02 Jan 06, 03:04:05 pm",
	}
	formats[5] = DateFormat{
		Display: "24hr",
		format:  "02 Jan 06, 15:04:05 pm",
	}
	return formats
}

func (a *App) getDateFormat() string {
	return getDateFormatData()[a.Settings.DateFormat].format
}

func (a *App) GetDateFormats() map[uint8]DateFormat {
	formats := getDateFormatData()
	now := time.Now()
	for i, f := range formats {
		f.Display = fmt.Sprintf("%s (%s)", now.Format(f.format), f.Display)
		formats[i] = f
	}

	return formats
}
