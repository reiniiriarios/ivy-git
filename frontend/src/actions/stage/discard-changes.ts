import { parseResponse } from "scripts/parse-response";
import { changes } from "stores/changes";
import { currentDiff } from "stores/diffs";
import { messageDialog } from "stores/message-dialog";
import { DiscardChanges } from "wailsjs/go/main/App";
import { get } from "svelte/store";

export function discardChanges(file: string) {
  let filename = file.replaceAll('\\','/').substring(file.lastIndexOf('/') + 1);
  messageDialog.confirm({
    heading: "Discard Changes",
    message: `Are you sure you want to discard all changes in <strong>${filename}</strong>?`,
    confirm: "Discard",
    callbackConfirm: () => {
      DiscardChanges(file).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          if (file == get(currentDiff).File) {
            currentDiff.refresh();
          }
        });
      });
    }
  })
}
