import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { inProgressCommitMessage } from "stores/ui";
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
        inProgressCommitMessage.fetch();
      });
    },
  });
}

export default squashMergeBranch;
