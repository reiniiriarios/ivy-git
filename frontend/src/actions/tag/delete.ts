import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { DeleteTag } from "wailsjs/go/main/App";

function deleteTag(tag: string) {
  messageDialog.confirm({
    heading: 'Delete Tag',
    message: `Are you sure you want to delete the tag <strong>${tag}</strong>?`,
    confirm: 'Delete',
    okay: 'Cancel',
    callbackConfirm: () => {
      DeleteTag(tag).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    }
  });
}

export default deleteTag;
