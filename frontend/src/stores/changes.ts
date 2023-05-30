import { parseResponse } from 'scripts/parse-response';
import { get, writable } from 'svelte/store';
import { GetUnstagedFileParsedDiff, GitListChanges } from 'wailsjs/go/main/App';
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
  Lines: ChangeLine[];
}

interface ChangeLine {
  RawNo: number;
  OldNo: number;
  NewNo: number;
  Selected: boolean;
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
      if (!f) return;
      if (xy === 'y') {
        GetUnstagedFileParsedDiff(file, f.Letter).then(result => {
          parseResponse(result, () => {
            update(c => {
              c[xy][file].Diff = result.Data;
              return c;
            });
            result.Data.File = file;
            result.Data.Status = f.Letter;
            result.Data.Staged = false;
            result.Data.Committed = false;
            currentDiff.set(result.Data);
          });
        });
      }
      else {
        //...
      }
    }
  };
}
export const changes = createChanges();
