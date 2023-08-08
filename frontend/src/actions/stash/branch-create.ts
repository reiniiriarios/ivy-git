import { checkRef } from "scripts/check-ref";
import { parseResponse } from "scripts/parse-response";
import { currentBranch } from "stores/branches";
import { commitData, commitSignData } from "stores/commits";
import { currentCommit } from "stores/commit-details";
import { messageDialog } from "stores/message-dialog";
import { CreateBranchFromStash } from "wailsjs/go/ivy/App";

function createBranchFromStash(stash: string, hash: string) {
  messageDialog.confirm({
    heading: 'Create Branch from Stash',
    message: `Create a branch from the stash <strong>${stash}</strong>:`,
    blank: 'Name of Branch',
    validateBlank: checkRef,
    confirm: 'Create',
    callbackConfirm: () => {
      CreateBranchFromStash(stash, messageDialog.blankValue()).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
          currentCommit.clearIfCurrent(hash);
          currentBranch.refresh();
        })
      });
    }
  });
}

export default createBranchFromStash;
