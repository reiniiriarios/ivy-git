import { parseResponse } from "scripts/parse-response";
import { messageDialog } from "stores/message-dialog";
import { remoteData } from "stores/remotes";
import { DeleteRemote } from "wailsjs/go/main/App";

function deleteRemote(remote: string) {
  messageDialog.confirm({
    heading: 'Delete Remote',
    message: `Are you sure you want to delete the remote <strong>${remote}</strong>?`,
    confirm: 'Delete',
    okay: 'Cancel',
    callbackConfirm: () => {
      DeleteRemote(remote).then(r => {
        parseResponse(r, () => {
          remoteData.refresh();
        });
      });
    }
  });
}

export default deleteRemote;
