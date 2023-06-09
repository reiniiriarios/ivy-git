import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { commitMessage } from "stores/ui";
import { CherryPick } from "wailsjs/go/main/App";

function cherryPick(hash: string) {
  messageDialog.confirm({
    heading: 'Cherry Pick Commit',
    message: `Cherry pick commit <strong>${hash.substring(0, 7)}</strong>.`,
    checkboxes: [
      {
        id: 'record',
        label: 'Record Original Hash',
      },
      {
        id: 'no_commit',
        label: 'No Commit',
      },
    ],
    confirm: 'Cherry Pick',
    callbackConfirm: () => {
      let no_commit = messageDialog.tickboxTicked('no_commit');
      CherryPick(hash, messageDialog.tickboxTicked('record'), no_commit).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
          if (no_commit) {
            commitMessage.fetch();
          }
        }, () => {
          commitMessage.fetch();
        });
      });
    }
  });
}

export default cherryPick;
