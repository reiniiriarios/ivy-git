import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { HardReset } from "wailsjs/go/main/App";

function hardReset(hash: string) {
  messageDialog.confirm({
    heading: 'Hard Reset',
    message: `Are you sure you want to hard reset to <strong>${hash.substring(0, 7)}</strong>?`,
    confirm: 'Hard Reset',
    callbackConfirm: () => {
      HardReset(hash).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default hardReset;
