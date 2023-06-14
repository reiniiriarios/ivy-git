import { branches } from 'stores/branches';
import { changes } from 'stores/changes';
import { commitData, commitSignData } from 'stores/commits';
import { currentDiff } from 'stores/diffs';
import { remoteData } from 'stores/remotes';
import { repoState } from 'stores/repo-state';
import { EventsOn } from 'wailsjs/runtime/runtime';

interface WatcherEvent {
	CommitChange: boolean;
	ShowRefChange: boolean;
	UncommittedDiffChange: boolean;
  UntrackedFilesChange: boolean;
	RemoteDiffChange: boolean;
	StagedDiffChange: boolean;
}

export function enableWatcher() {
  EventsOn('watcher', (e: WatcherEvent) => {
    console.log('Watcher updating...');
    repoState.refresh();
    changes.refresh();
    commitData.refresh();
    currentDiff.refresh();
    if (e.CommitChange || e.ShowRefChange || e.UncommittedDiffChange || e.RemoteDiffChange) {
      branches.refresh();
      commitSignData.refresh();
      remoteData.refresh();
    }
  });
}
