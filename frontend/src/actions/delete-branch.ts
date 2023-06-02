import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { DeleteBranch } from "wailsjs/go/main/App";

// Delete local branch, optionally remote.
function deleteBranch(branch: string, can_delete_remote: boolean) {
  let opts = [{id: 'force', label: 'Force Delete'}];
  if (can_delete_remote) {
    opts.push({id: 'remote', label: 'Delete on Remote'});
  }
  messageDialog.confirm({
    heading: 'Delete Branch',
    message: `Are you sure you want to delete the branch <strong>${branch}</strong>?`,
    confirm: 'Delete',
    okay: 'Cancel',
    checkboxes: opts,
    callbackConfirm: () => {
      let remote = messageDialog.tickboxTicked('remote');
      DeleteBranch(branch, messageDialog.tickboxTicked('force'), remote).then(r => {
        if (r.Response === 'must-force') {
          messageDialog.confirm({
            heading: 'Force Delete Branch',
            message: `The branch <strong>${branch}</strong> could not be deleted because it is not fully merged.\n\nWould you like to force delete the branch?`,
            confirm: 'Force Delete',
            okay: 'Cancel',
            callbackConfirm: () => {
              DeleteBranch(branch, true, remote).then(r => {
                parseResponse(r, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                });
              });
            }
          });
        } else {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        }
      });
    },
  });
}

export default deleteBranch;
