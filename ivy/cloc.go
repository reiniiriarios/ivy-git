package ivy

import (
	"errors"
	"fmt"
	"ivy-git/ivy/git"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

// Update cloc data with latest git information.
func (a *App) updateClocData(clocData git.ClocData) git.ClocData {
	newClocData, err := a.Git.Cloc()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	} else {
		a.saveClocData(newClocData)
		clocData = newClocData
	}
	return clocData
}

// Save cloc data to config file.
func (a *App) saveClocData(clocData git.ClocData) {
	data, err := yaml.Marshal(&clocData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	err = os.WriteFile(a.clocFileLocation(), []byte(data), 0644)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
}

// Load existing loc data, if any.
func (a *App) loadClocData() git.ClocData {
	if a.RepoSaveData.CurrentRepo == "" {
		// This should never happen, as a repo should need to be selected before this command can be run.
		runtime.LogError(a.ctx, "no current repo selected")
		return git.ClocData{}
	}

	raw := a.initClocData()
	var clocData git.ClocData
	err := yaml.Unmarshal(raw, &clocData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	return clocData
}

// Return file data, create empty file if not found.
func (a *App) initClocData() []byte {
	file := a.clocFileLocation()
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, e := os.Create(file)
		if e != nil {
			runtime.LogError(a.ctx, e.Error())
		}
		defer f.Close()
	}
	contents, err := os.ReadFile(file)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	return contents
}

func (a *App) resetClocData() {
	file := a.clocFileLocation()
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, e := os.Create(file)
		if e != nil {
			runtime.LogError(a.ctx, e.Error())
		}
		defer f.Close()
	} else {
		os.Truncate(file, 0)
	}
}

// Location of cloc file for current repo.
func (a *App) clocFileLocation() string {
	return filepath.Join(a.settingsLocationLocal(), fmt.Sprintf("cloc.%s.yaml", a.RepoSaveData.CurrentRepo))
}
