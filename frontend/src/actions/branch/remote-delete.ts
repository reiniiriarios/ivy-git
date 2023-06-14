import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { DeleteRemoteBranch } from "wailsjs/go/main/App";

// Delete only remote branch.
function deleteRemoteBranch(branch: string, remote: string) {
  messageDialog.confirm({
    heading: 'Delete Remote Branch',
    message: `Are you sure you want to delete the remote branch <strong>${remote}/${branch}</strong>?`,
    confirm: 'Delete',
    okay: 'Cancel',
    checkboxes: [{id: 'force', label: 'Force Delete'}],
    callbackConfirm: () => {
      DeleteRemoteBranch(branch, remote, messageDialog.tickboxTicked('force')).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default deleteRemoteBranch;
