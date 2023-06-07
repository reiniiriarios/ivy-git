import { parseResponse } from 'scripts/parse-response';
import { derived, get, writable } from 'svelte/store';
import { GetConflictParsedDiff, GetWorkingFileParsedDiff, GitListChanges } from 'wailsjs/go/main/App';
import { currentRepo } from 'stores/repos';
import { currentDiff, type Diff } from 'stores/diffs';
import { currentTab } from 'stores/ui';

interface Changes {
  x: Change[], // staged
  y: Change[], // unstaged
  c: Change[], // conflicts
}

interface Change {
  File: string;
  Basename: string;
  Dir: string;
  Letter: string;
  Them: string;
  Us: string;
  Flag: string;
  // Addl Fetch
  Diff: Diff;
  // UI
  Partial: boolean;
  Resolved: boolean;
}

function createChanges() {
  const { subscribe, set, update } = writable({x: [], y: []} as Changes);

  return {
    subscribe,
    refresh: async () => {
      if (get(currentRepo)) {
        GitListChanges().then(result => {
          parseResponse(result, () => {
            set({
              x: result.Data.ChangesX ?? [],
              y: result.Data.ChangesY ?? [],
              c: result.Data.ChangesC ?? [],
            });
            currentDiff.refresh();
            // If there's nothing left to do, switch away from the changes tab.
            if (!changes.numStaged() && !changes.numUnstaged() && !changes.numConflicts()) {
              if (get(currentTab) === 'changes') {
                currentTab.set('tree');
              }
            }
          });
        });
      } else {
        set({x: [], y: [], c: []});
      }
    },
    fetchDiff: async (xyc: string, file: string, ignoreCache: boolean = false) => {
      // The file will be in the X (staged), Y (unstaged), or C (conflicts) list.
      if (xyc !== 'x' && xyc !== 'c') xyc = 'y';
      let f = get(changes)[xyc][file];
      // If file not in changes list, clear current diff as it's outdated.
      if (!f) {
        currentDiff.clear();
        return;
      }
      // If the diff is already fetched, display that, don't refetch. Unless we say so.
      if (f.Diff && !ignoreCache) {
        currentDiff.set(f.Diff);
      }
      else {
        // Conflict diff.
        if (xyc === 'c') {
          GetConflictParsedDiff(file).then(result => {
            parseResponse(result, () => {
              result.Data.File = file;
              result.Data.Status = f.Letter;
              result.Data.Conflict = true;
              result.Data.Staged = false;
              result.Data.Committed = false;
              result.Data.ConflictSelections = [];
              update(c => {
                if (c.c[file]) {
                  c.c[file].Diff = result.Data;
                }
                return c;
              });
              currentDiff.set(result.Data);
            });
          });
        }
        // Staged or unstaged diff.
        else {
          GetWorkingFileParsedDiff(file, f.Letter, xyc === 'x').then(result => {
            parseResponse(result, () => {
              result.Data.File = file;
              result.Data.Status = f.Letter;
              result.Data.Staged = xyc === 'x';
              result.Data.Committed = false;
              update(c => {
                if (c[xyc][file]) {
                  c[xyc][file].Diff = result.Data;
                }
                return c;
              });
              currentDiff.set(result.Data);
            });
          });
        }
      }
    },
    setPartial: async (xy: string, file: string, partial: boolean) => {
      if (xy !== 'x') xy = 'y';
      update(c => {
        if (c[xy][file]) {
          c[xy][file].Partial = partial;
        }
        return c;
      });
    },
    setResolved: async (file: string, resolved: boolean) => {
      update(c => {
        if (c.c[file]) {
          c.c[file].Resolved = resolved;
        }
        return c;
      })
    },
    numStaged: () => {
      let c = get(changes);
      return c?.x ? Object.keys(get(changes).x).length : 0;
    },
    numUnstaged: () => {
      let c = get(changes);
      return c?.y ? Object.keys(get(changes).y).length : 0;
    },
    numConflicts: () => {
      let c = get(changes);
      return c?.c ? Object.keys(get(changes).c).length : 0;
    },
  };
}
export const changes = createChanges();
export const mergeConflicts = derived(changes, $changes => $changes?.c && Object.keys($changes.c).length);
export const mergeConflictsResolved = derived(changes, $changes => $changes?.c && Object.values($changes.c).every(c => c.Resolved));
