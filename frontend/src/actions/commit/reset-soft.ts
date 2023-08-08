import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { SoftReset } from "wailsjs/go/ivy/App";

function softReset(hash: string) {
  messageDialog.confirm({
    heading: 'Soft Reset',
    message: `Are you sure you want to soft reset to <strong>${hash.substring(0, 7)}</strong>?`,
    confirm: 'Soft Reset',
    callbackConfirm: () => {
      SoftReset(hash).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default softReset;
