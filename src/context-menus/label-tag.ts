import { ClipboardSetText } from "src/_tmp";
import { type Menu } from "context-menus/_all";
import { commitData, commitSignData } from "stores/commit-data";
import { PushTag, DeleteTag } from "src/_tmp";
import { parseResponse } from "scripts/parse-response";
import { messageDialog } from "stores/message-dialog";

export const menuLabelTag: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Push Tag",
      callback: () => {
        PushTag(e.dataset.name).then(r => {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        });
      },
    },
    {
      text: "Delete Tag",
      callback: () => {
        messageDialog.confirm({
          heading: 'Delete Tag',
          message: `Are you sure you want to delete the tag <strong>${e.dataset.name}</strong>?`,
          confirm: 'Delete',
          okay: 'Cancel',
          callbackConfirm: () => {
            DeleteTag(e.dataset.name).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          }
        });
      },
    },
    {
      sep: true,
    },
    {
      text: "Copy Tag Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
  ];
}
