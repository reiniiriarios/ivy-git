import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { commitMessage } from "stores/ui";
import { MergeRebase } from "wailsjs/go/main/App";

function rebaseAndMergeIntoCurrentBranch(branch: string) {
  messageDialog.confirm({
    heading: 'Rebase and Merge into Current Branch',
    message: `Rebase <strong>${branch}</strong> onto current branch and merge?`,
    confirm: 'Merge',
    okay: 'Cancel',
    callbackConfirm: () => {
      MergeRebase(branch).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        }, () => {
          commitMessage.fetch();
        });
      });
    },
  });
}

export default rebaseAndMergeIntoCurrentBranch;
