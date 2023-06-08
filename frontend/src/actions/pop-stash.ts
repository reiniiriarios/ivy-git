import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { currentCommit } from "stores/commit-details";
import { messageDialog } from "stores/message-dialog";
import { PopStash } from "wailsjs/go/main/App";

function popStash(stash: string, hash: string) {
  messageDialog.confirm({
    heading: 'Pop Stash',
    message: `Pop the stash <strong>${stash}</strong> onto the current working tree?`,
    confirm: 'Pop',
    callbackConfirm: () => {
      PopStash(stash, messageDialog.tickboxTicked('index')).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
          currentCommit.clearIfCurrent(hash);
        })
      });
    }
  });
}

export default popStash;
