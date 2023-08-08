package ivy

import (
	"errors"
	"fmt"
	"ivy-git/ivy/git"
	"os"
	"path/filepath"
	"sort"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

type ContributorData struct {
	LastHashParsed string
	Contributors   git.Contributors
}

// Update contributor data with latest git information.
func (a *App) updateContributorData(contributorData ContributorData) ContributorData {
	var err error
	contributorData.Contributors, contributorData.LastHashParsed, err = a.Git.AddContributorsSince(contributorData.Contributors, contributorData.LastHashParsed)
	sort.Sort(contributorData.Contributors)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	} else {
		a.saveContributorData(contributorData)
	}
	return contributorData
}

// Save contributor data to config file.
func (a *App) saveContributorData(contributorData ContributorData) {
	data, err := yaml.Marshal(&contributorData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	err2 := os.WriteFile(a.contributorsFileLocation(), []byte(data), 0644)
	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}
}

// Load existing contributor data, if any.
func (a *App) loadContributorData() ContributorData {
	if a.RepoSaveData.CurrentRepo == "" {
		// This should never happen, as a repo should need to be selected before this command can be run.
		runtime.LogError(a.ctx, "no current repo selected")
		return ContributorData{}
	}

	raw := a.initContributorData()
	var contributorData ContributorData
	err := yaml.Unmarshal(raw, &contributorData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
	return contributorData
}

// Return file data, create empty file if not found.
func (a *App) initContributorData() []byte {
	file := a.contributorsFileLocation()
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

func (a *App) resetContributorData() {
	file := a.contributorsFileLocation()
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

// Location of contributors file for current repo.
func (a *App) contributorsFileLocation() string {
	return filepath.Join(a.settingsLocationLocal(), fmt.Sprintf("contributors.%s.yaml", a.RepoSaveData.CurrentRepo))
}
