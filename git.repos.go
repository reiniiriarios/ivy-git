package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

type RepoSaveData struct {
	CurrentRepo string
	Repos       map[string]Repo
}

type Repo struct {
	Name      string
	Directory string
}

type RepoResponse struct {
	Response string
	Message  string
	Id       string
	Repo     Repo
}

func (a *App) LoadYaml() {
	if _, err := os.Stat(a.repoYamlLocation()); errors.Is(err, os.ErrNotExist) {
		a.initYaml()
	}

	yfile, err := ioutil.ReadFile(a.repoYamlLocation())

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	var data RepoSaveData

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}

	a.RepoSaveData = data
}

func (a *App) GetRepos() map[string]Repo {
	return a.RepoSaveData.Repos
}

func (a *App) UpdateSelectedRepo(repo string) {
	a.RepoSaveData.CurrentRepo = repo
	a.SaveRepoData()
}

func (a *App) GetSelectedRepo() string {
	return a.RepoSaveData.CurrentRepo
}

func (a *App) AddRepo() RepoResponse {
	d, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Git Repository",
	})

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return RepoResponse{
			Response: "error",
			Message:  err.Error(),
		}
	}

	if d == "" {
		return RepoResponse{
			Response: "",
		}
	}

	if !a.IsGitRepo(d) {
		return RepoResponse{
			Response: "error",
			Message:  "The directory selected is not a git repository.",
		}
	}

	for _, r := range a.RepoSaveData.Repos {
		if r.Directory == d {
			return RepoResponse{
				Response: "error",
				Message:  "The directory selected is already added.",
			}
		}
	}

	id := uuid.New().String()
	newRepo := Repo{
		Name:      filepath.Base(d),
		Directory: d,
	}

	if a.RepoSaveData.Repos == nil {
		a.RepoSaveData.Repos = make(map[string]Repo)
	}
	a.RepoSaveData.Repos[id] = newRepo
	a.SaveRepoData()

	return RepoResponse{
		Response: "success",
		Id:       id,
		Repo:     newRepo,
	}
}

func (a *App) RemoveRepo(id string) map[string]Repo {
	delete(a.RepoSaveData.Repos, id)
	a.SaveRepoData()
	return a.RepoSaveData.Repos
}

func (a *App) SaveRepoData() {
	data, err := yaml.Marshal(&a.RepoSaveData)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	err2 := os.WriteFile(a.repoYamlLocation(), []byte(data), 0644)
	if err2 != nil {
		runtime.LogError(a.ctx, err2.Error())
	}
}

func (a *App) initYaml() {
	f, e := os.Create(a.repoYamlLocation())
	if e != nil {
		runtime.LogError(a.ctx, e.Error())
	}
	defer f.Close()
}

func (a *App) repoYamlLocation() string {
	// todo
	return "repos.yaml"
}
