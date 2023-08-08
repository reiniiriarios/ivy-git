import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { AddTag } from "wailsjs/go/ivy/App";

function addTag(hash: string) {
  messageDialog.addTag({
    message: `Add tag to commit <strong>${hash.substring(0, 7)}</strong>:`,
    callbackConfirm: () => {
      let data = messageDialog.addTagData();
      AddTag(hash, data.name, data.message, data.push).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default addTag;
