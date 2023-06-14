import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { commitMessage } from "stores/ui";
import { MergeSquash } from "wailsjs/go/main/App";

function squashMergeBranch(branch: string) {
  messageDialog.confirm({
    heading: 'Squash & Merge onto Current Branch',
    message: `Squash <strong>${branch}</strong> and merge onto current branch?`,
    confirm: 'Merge',
    okay: 'Cancel',
    callbackConfirm: () => {
      MergeSquash(branch).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
        commitMessage.fetch();
      });
    },
  });
}

export default squashMergeBranch;
