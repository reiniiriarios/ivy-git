import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { RevertCommit } from "wailsjs/go/main/App";

function revertCommit(hash: string) {
  messageDialog.confirm({
    heading: 'Revert Commit',
    message: `Are you sure you want to revert <strong>${hash.substring(0, 7)}</strong>?`,
    confirm: 'Revert',
    callbackConfirm: () => {
      RevertCommit(hash).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default revertCommit;
