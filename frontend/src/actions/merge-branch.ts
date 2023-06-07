import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { commitMessage } from "stores/ui";
import { MergeCommit } from "wailsjs/go/main/App";

function mergeBranch(branch: string) {
  messageDialog.confirm({
    heading: 'Merge into Current Branch',
    message: `Merge <strong>${branch}</strong> into current branch?`,
    confirm: 'Merge',
    okay: 'Cancel',
    checkboxes: [
      {
        id: 'no-ff',
        label: 'Create a new commit even if fast-forward is possible',
        checked: true,
      },
      {
        id: 'no-commit',
        label: 'No Commit',
        checked: false,
      },
    ],
    callbackConfirm: () => {
      let no_commit = messageDialog.tickboxTicked('no-commit');
      let no_ff = messageDialog.tickboxTicked('no-ff');
      MergeCommit(branch, no_commit, no_ff).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
          if (no_commit) {
            commitMessage.fetch();
          }
        }, () => {
          commitMessage.fetch();
        });
      });
    },
  });
}

export default mergeBranch;
