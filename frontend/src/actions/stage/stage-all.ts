import { get } from "svelte/store";
import { parseResponse } from "scripts/parse-response";
import { changes } from "stores/changes";
import { currentDiff } from "stores/diffs";
import { currentTab } from "stores/ui";
import { StageAll } from "wailsjs/go/main/App";

export function stageAll() {
  StageAll().then(result => {
    parseResponse(result, () => {
      changes.refresh();
      if (get(currentTab) === 'changes') {
        if (get(currentDiff).Staged) {
          currentDiff.refresh();
        } else {
          currentDiff.clear();
        }
      }
    });
  });
}
