import type { Menu } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { messageDialog } from "stores/message-dialog";
import { remoteData } from "stores/remotes";
import { DeleteRemote } from "wailsjs/go/main/App";

export const menuRemote: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Delete Remote",
      callback: () => {
        messageDialog.confirm({
          heading: 'Delete Remote',
          message: `Are you sure you want to delete the remote <strong>${e.dataset.name}</strong>?`,
          confirm: 'Delete',
          okay: 'Cancel',
          callbackConfirm: () => {
            DeleteRemote(e.dataset.name).then(r => {
              parseResponse(r, () => {
                remoteData.refresh();
              });
            });
          }
        });
      },
    },
  ];
}
