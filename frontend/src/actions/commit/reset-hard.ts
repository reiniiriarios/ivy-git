import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { HardReset } from "wailsjs/go/ivy/App";

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
