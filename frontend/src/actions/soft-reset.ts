import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { SoftReset } from "wailsjs/go/main/App";

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
