import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { ApplyStash } from "wailsjs/go/main/App";

function applyStash(stash: string) {
  messageDialog.confirm({
    heading: 'Apply Stash',
    message: `Apply the stash <strong>${stash}</strong> onto the current working tree?`,
    confirm: 'Apply',
    checkboxes: [{
      id: 'index',
      label: 'Reinstate Index',
      checked: false,
    }],
    callbackConfirm: () => {
      ApplyStash(stash, messageDialog.tickboxTicked('index')).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        })
      });
    }
  });
}

export default applyStash;
