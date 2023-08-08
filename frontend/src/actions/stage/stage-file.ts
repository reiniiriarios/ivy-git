import { parseResponse } from "scripts/parse-response";
import { changes, type Change } from "stores/changes";
import { currentDiff } from "stores/diffs";
import { currentTab } from "stores/ui";
import { get } from "svelte/store";
import { StageFile, StagePartialFile } from "wailsjs/go/ivy/App";

export function stageFile(file: Change, partial: boolean) {
  if (file) {
    if (partial) {
      StagePartialFile(file.Diff, file.File, file.Letter).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          let diff = get(currentDiff);
          if (diff.File === file.File && !diff.Staged && get(currentTab) === 'changes') {
            currentDiff.refresh();
          }
        });
      });
    } else {
      StageFile(file.File).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          let diff = get(currentDiff);
          if (diff.File === file.File && !diff.Staged && get(currentTab) === 'changes') {
            currentDiff.clear();
          }
        });
      });
    }
  }
}
