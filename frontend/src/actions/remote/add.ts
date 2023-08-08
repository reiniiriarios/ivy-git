import { parseResponse } from "scripts/parse-response";
import { messageDialog } from "stores/message-dialog";
import { remoteData } from "stores/remotes";
import { AddRemote } from "wailsjs/go/ivy/App";

function addRemote() {
  messageDialog.addRemote({
    callbackConfirm: () => {
      let data = messageDialog.addRemoteData();
      AddRemote(data.name, data.fetch, data.push).then(result => {
        parseResponse(result, () => {
          remoteData.refresh();
          messageDialog.clear();
        }, () => {
          // Refresh regardless of error. The remote may be added, but
          // an error may occur when trying to fetch from it. The list
          // should still be updated.
          remoteData.refresh();
        });
      });
    },
  });
}

export default addRemote;
