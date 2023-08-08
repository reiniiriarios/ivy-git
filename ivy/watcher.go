package ivy

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const WATCHER_INTERVAL = 2

type WatcherEvent struct {
	CommitChange          bool
	ShowRefChange         bool
	UncommittedDiffChange bool
	UntrackedFilesChange  bool
	RemoteDiffChange      bool
	StagedDiffChange      bool
}

// Watcher for changes in repo.
func (a *App) watcher() {
	// Increment a (sort of but not really) semaphore to ensure
	// that only one instance of this loop is run at a time.
	a.WatcherSemiSemaphore++
	semi_semaphore := a.WatcherSemiSemaphore

	var wg sync.WaitGroup

	// Set initial data.
	if a.isCurrentRepo() {
		wg.Add(6)
		go a.updateLastCommit(nil, &wg)
		go a.updateShowRefAll(nil, &wg)
		go a.updateUncommittedDiff(nil, &wg)
		go a.updateUntrackedFiles(nil, &wg)
		go a.updateRemoteDiff(nil, &wg)
		go a.updateStagedDiff(nil, &wg)
		wg.Wait()
	}

	for range time.Tick(time.Second * WATCHER_INTERVAL) {
		// If this variable has changed, it means another instance of this
		// loop is running and this one should quit.
		if semi_semaphore != a.WatcherSemiSemaphore {
			return
		}

		// If no repo is selected, don't do anything.
		if !a.isCurrentRepo() {
			continue
		}

		// Get new data
		wg.Add(6)
		var lc_new bool
		var sr_new bool
		var ud_new bool
		var uf_new bool
		var rd_new bool
		var sd_new bool
		go a.updateLastCommit(&lc_new, &wg)
		go a.updateShowRefAll(&sr_new, &wg)
		go a.updateUncommittedDiff(&ud_new, &wg)
		go a.updateUntrackedFiles(&uf_new, &wg)
		go a.updateRemoteDiff(&rd_new, &wg)
		go a.updateStagedDiff(&sd_new, &wg)
		wg.Wait()

		// If data changed, emit event.
		if lc_new || sr_new || ud_new || uf_new || rd_new || sd_new {
			runtime.LogInfo(a.ctx, "Watcher updating")
			runtime.EventsEmit(a.ctx, "watcher", WatcherEvent{
				CommitChange:          lc_new,
				ShowRefChange:         sr_new,
				UncommittedDiffChange: ud_new,
				UntrackedFilesChange:  uf_new,
				RemoteDiffChange:      rd_new,
				StagedDiffChange:      sd_new,
			})
		}
	}
}

// Update last commit hash for watcher.
func (a *App) updateLastCommit(new *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	last_commit, err := a.Git.GetLastCommitHash()
	if err != nil && new != nil {
		*new = true
	} else if new != nil {
		*new = a.CurrentHash != last_commit
	}
	a.CurrentHash = last_commit
}

// Update md5 of show refs for watcher.
func (a *App) updateShowRefAll(new *bool, wg *sync.WaitGroup) {
	diff, err := a.Git.ShowRefAll()
	a.updateWatcherDiff(new, wg, &a.ShowRefAll, diff, err)
}

// Update list of untracked files for watcher.
func (a *App) updateUntrackedFiles(new *bool, wg *sync.WaitGroup) {
	diff, err := a.Git.GetUntrackedFiles()
	a.updateWatcherDiff(new, wg, &a.UntrackedFiles, diff, err)
}

// Update md5 of uncommitted diff for watcher.
func (a *App) updateUncommittedDiff(new *bool, wg *sync.WaitGroup) {
	diff, err := a.Git.GetUncommittedDiff()
	a.updateWatcherDiff(new, wg, &a.UncommittedDiff, diff, err)
}

// Update md5 of remote diff for watcher.
func (a *App) updateRemoteDiff(new *bool, wg *sync.WaitGroup) {
	// This will error when the local branch doesn't have a remote, so we ignore.
	diff, _ := a.Git.GetDiffRemoteCurrent()
	a.updateWatcherDiff(new, wg, &a.RemoteDiff, diff, nil)
}

// Update md5 of staged diff for watcher.
func (a *App) updateStagedDiff(new *bool, wg *sync.WaitGroup) {
	diff, err := a.Git.GetDiffStaged()
	a.updateWatcherDiff(new, wg, &a.StagedDiff, diff, err)
}

// Update hash of diff for watcher.
func (a *App) updateWatcherDiff(new *bool, wg *sync.WaitGroup, update_variable *string, diff string, err error) {
	defer wg.Done()

	if err != nil && new != nil {
		*new = true
	} else {
		hash := md5.Sum([]byte(diff))
		hash_hex := hex.EncodeToString(hash[:])
		if new != nil {
			*new = hash_hex != *update_variable
		}
		*update_variable = hash_hex
	}
}
