package main

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WatcherEvent struct {
	CommitChange          bool
	UncommittedDiffChange bool
	RemoteDiffChange      bool
	StagedDiffChange      bool
}

// Watcher for changes in repo.
func (a *App) watcher() {
	// Increment a (sort of but not really) semaphore to ensure
	// that only one instance of this loop is run at a time.
	a.WatcherSemiSemaphore++
	semi_semaphore := a.WatcherSemiSemaphore

	for range time.Tick(time.Second * 2) {
		// If this variable has changed, it means another instance of this
		// loop is running and this one should quit.
		if semi_semaphore != a.WatcherSemiSemaphore {
			return
		}

		lc_new := a.updateLastCommit()
		ud_new := a.updateUncommittedDiff()
		rd_new := a.updateRemoteDiff()
		sd_new := a.updateStagedDiff()
		if lc_new || ud_new || rd_new || sd_new {
			runtime.LogInfo(a.ctx, "Watcher updating")
			runtime.EventsEmit(a.ctx, "watcher", WatcherEvent{
				CommitChange:          lc_new,
				UncommittedDiffChange: ud_new,
				RemoteDiffChange:      rd_new,
				StagedDiffChange:      sd_new,
			})
		}
	}
}

// Update last commit hash for watcher.
func (a *App) updateLastCommit() bool {
	last_commit, err := a.Git.GetLastCommitHash()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return true
	}
	new := a.CurrentHash != last_commit
	a.CurrentHash = last_commit
	return new
}

// Update md5 of uncommitted diff for watcher.
func (a *App) updateUncommittedDiff() bool {
	diff, err := a.Git.GetUncommittedDiff()
	return a.updateWatcherDiff(&a.UncommittedDiff, diff, err)
}

// Update md5 of remote diff for watcher.
func (a *App) updateRemoteDiff() bool {
	diff, err := a.Git.GetDiffRemoteCurrent()
	return a.updateWatcherDiff(&a.RemoteDiff, diff, err)
}

// Update md5 of staged diff for watcher.
func (a *App) updateStagedDiff() bool {
	diff, err := a.Git.GetDiffStaged()
	return a.updateWatcherDiff(&a.StagedDiff, diff, err)
}

// Update hash of diff for watcher.
func (a *App) updateWatcherDiff(update_variable *string, diff string, err error) bool {
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return true
	}
	hash := md5.Sum([]byte(diff))
	hash_hex := hex.EncodeToString(hash[:])
	new := hash_hex != *update_variable
	*update_variable = hash_hex
	return new
}
