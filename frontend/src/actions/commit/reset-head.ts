import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { ResetHead } from "wailsjs/go/ivy/App";

function resetHead(hard: boolean) {
  let word = hard ? 'Hard' : 'Soft';
  messageDialog.confirm({
    heading: `${word} Reset`,
    message: `Are you sure you want to ${word.toLowerCase()} reset the last commit?`,
    confirm: `${word} Reset`,
    callbackConfirm: () => {
      ResetHead(hard).then(result => {
        parseResponse(result, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default resetHead;
