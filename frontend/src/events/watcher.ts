import { branches } from 'stores/branches';
import { changes } from 'stores/changes';
import { commitData, commitSignData } from 'stores/commit-data';
import { currentDiff } from 'stores/diffs';
import { remoteData } from 'stores/remotes';
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
