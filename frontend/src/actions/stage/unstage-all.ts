import { get } from "svelte/store";
import { parseResponse } from "scripts/parse-response";
import { changes } from "stores/changes";
import { currentDiff } from "stores/diffs";
import { currentTab } from "stores/ui";
import { UnstageAll } from "wailsjs/go/ivy/App";

export function unstageAll() {
  UnstageAll().then(result => {
    parseResponse(result, () => {
      changes.refresh();
      if (get(currentTab) === 'changes') {
        if (get(currentDiff).Staged) {
          currentDiff.clear();
        } else {
          currentDiff.refresh();
        }
      }
    });
  });
}
