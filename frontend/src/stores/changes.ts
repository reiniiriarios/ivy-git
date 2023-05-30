import { parseResponse } from 'scripts/parse-response';
import { get, writable } from 'svelte/store';
import { GetWorkingFileParsedDiff, GitListChanges } from 'wailsjs/go/main/App';
import { currentRepo } from 'stores/repos';
import { currentDiff, type Diff } from 'stores/diffs';

interface Changes {
  x: Change[],
  y: Change[],
}

interface Change {
  File: string;
  Basename: string;
  Dir: string;
  Letter: string;
  Flag: string;
  // Addl Fetch
  Diff: Diff;
  // UI
  Partial: boolean;
}

function createChanges() {
  const { subscribe, set, update } = writable({x: [], y: []} as Changes);

  return {
    subscribe,
    refresh: async () => {
      if (get(currentRepo)) {
        GitListChanges().then(result => {
          parseResponse(result, () => set({
            x: result.Data.ChangesX ?? [],
            y: result.Data.ChangesY ?? [],
          }));
        });
      } else {
        set({x: [], y: []});
      }
    },
    fetchDiff: async (xy: string, file: string) => {
      if (xy !== 'x') xy = 'y';
      let f = get(changes)[xy][file];
      // If file not in changes list, clear current diff as it's outdated.
      if (!f) {
        currentDiff.clear();
        return;
      }
      if (f.Diff) {
        currentDiff.set(f.Diff);
      } else {
        GetWorkingFileParsedDiff(file, f.Letter, xy === 'x').then(result => {
          parseResponse(result, () => {
            update(c => {
              if (c[xy][file]) {
                c[xy][file].Diff = result.Data;
              }
              return c;
            });
            result.Data.File = file;
            result.Data.Status = f.Letter;
            result.Data.Staged = xy === 'x';
            result.Data.Committed = false;
            currentDiff.set(result.Data);
          });
        });
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
    numStaged: () => {
      return Object.keys(get(changes).x).length;
    },
    numUnstaged: () => {
      return Object.keys(get(changes).y).length;
    },
  };
}
export const changes = createChanges();
