import { checkRef } from "scripts/check-ref";
import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { currentRepo, repos } from "stores/repos";
import { get } from "svelte/store";
import { RenameBranch } from "wailsjs/go/ivy/App";

function renameBranch(branch: string) {
  messageDialog.confirm({
    heading: 'Rename Branch',
    message: `Rename <strong>${branch}</strong> locally and on all remotes to:`,
    confirm: 'Rename',
    blank: 'New Name',
    validateBlank: checkRef,
    okay: 'Cancel',
    callbackConfirm: () => {
      let isMain = branch === get(repos)[get(currentRepo)].Main;
      RenameBranch(branch, messageDialog.blankValue()).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
          if (isMain) {
            repos.updateMain(messageDialog.blankValue());
          }
        });
      });
    },
  });
}

export default renameBranch;
