import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { settings } from "stores/settings";
import { get } from "svelte/store";
import { PushBranch } from "wailsjs/go/main/App";

function pushBranch(name: string, display: string) {
  PushBranch(name, false).then(r => {
    if (r.Response === 'must-force' && get(settings).Workflow !== 'merge') {
      messageDialog.confirm({
        heading: 'Force Push Branch',
        message: `Unable to push branch <strong>${display}</strong>, as it's behind its remote counterpart.\n\nForce push this branch?`,
        confirm: 'Force Push',
        okay: 'Cancel',
        callbackConfirm: () => {
          PushBranch(name, true).then(r => {
            parseResponse(r, () => {
              commitData.refresh();
              commitSignData.refresh();
            });
          });
        },
      });
    }
    else {
      parseResponse(r, () => {
        commitData.refresh();
        commitSignData.refresh();
      });
    }
  })
}

export default pushBranch;
