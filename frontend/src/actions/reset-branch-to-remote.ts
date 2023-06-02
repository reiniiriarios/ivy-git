import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { ResetBranchToRemote } from "wailsjs/go/main/App";

function resetBranchToRemote(branch: string) {
  messageDialog.confirm({
    heading: 'Reset Local Branch to Remote',
    message: `Are you sure you want to reset the local branch <strong>${branch}</strong> to its default remote?`,
    confirm: 'Reset',
    okay: 'Cancel',
    callbackConfirm: () => {
      ResetBranchToRemote(branch).then(() => {
        commitData.refresh();
        commitSignData.refresh();
      });
    },
  });
}

export default resetBranchToRemote;
