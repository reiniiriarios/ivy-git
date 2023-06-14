import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { currentCommit } from "stores/commit-details";
import { messageDialog } from "stores/message-dialog";
import { DropStash } from "wailsjs/go/main/App";

function dropStash(stash: string, hash: string) {
  messageDialog.confirm({
    heading: 'Drop Stash',
    message: `Are you sure you want to drop the stash <strong>${stash}</strong>?`,
    confirm: 'Drop',
    callbackConfirm: () => {
      DropStash(stash).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
          currentCommit.clearIfCurrent(hash);
        })
      });
    }
  });
}

export default dropStash;
