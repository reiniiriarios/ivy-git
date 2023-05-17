import { branches } from 'stores/branches';
import { changes } from 'stores/changes';
import { commitData, commitSignData } from 'stores/commit-data';
import { EventsOn } from 'wailsjs/runtime/runtime';

interface WatcherEvent {
	CommitChange: boolean;
	UncommittedDiffChange: boolean;
}

export function enableWatcher() {
  EventsOn('watcher', (e: WatcherEvent) => {
    console.log('Watcher updating...');
    branches.refresh();
    changes.refresh();
    commitData.refresh();
    commitSignData.refresh();
  });
}
