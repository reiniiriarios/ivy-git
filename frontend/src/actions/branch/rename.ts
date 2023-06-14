import { checkRef } from "scripts/check-ref";
import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { RenameBranch } from "wailsjs/go/main/App";

function renameBranch(branch: string) {
  messageDialog.confirm({
    heading: 'Rename Branch',
    message: `Rename <strong>${branch}</strong> locally and on all remotes to:`,
    confirm: 'Rename',
    blank: 'New Name',
    validateBlank: checkRef,
    okay: 'Cancel',
    callbackConfirm: () => {
      RenameBranch(branch, messageDialog.blankValue()).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default renameBranch;
