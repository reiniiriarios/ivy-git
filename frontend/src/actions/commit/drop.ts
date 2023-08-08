import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { DropCommit } from "wailsjs/go/ivy/App";

function dropCommit(hash: string) {
  messageDialog.confirm({
    heading: 'Drop Commit',
    message: `Are you sure you want to drop <strong>${hash.substring(0, 7)}</strong>?`,
    confirm: 'Drop',
    callbackConfirm: () => {
      DropCommit(hash).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default dropCommit;
