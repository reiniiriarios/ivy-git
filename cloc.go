package main

import "ivy-git/cloc"

func (a *App) Cloc() DataResponse {
	files, err := a.Git.LsFiles()
	if err != nil {
		return dataResponse(err, false)
	}
	cloc, err := cloc.Cloc(a.Git.Repo.Directory, files)
	return dataResponse(err, cloc)
}
