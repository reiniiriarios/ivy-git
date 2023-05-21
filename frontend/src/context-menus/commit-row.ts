import type { Menu, MenuItem } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { AddTag } from "wailsjs/go/main/App";

export const menuCommitRow: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.head !== 'true') {
    m = m.concat([
      {
        text: "Checkout Commit",
        callback: () => alert("todo: checkout"),
      },
      {
        text: "Cherry Pick Commit",
        callback: () => alert("todo: checkout"),
      },
    ]);
  }

  m = m.concat([
    {
      text: "Revert Commit",
      callback: () => alert("todo: revert"),
    },
    {
      text: "Add Tag",
      callback: () => {
        messageDialog.addTag({
          message: `Add tag to commit <strong>${e.dataset.hash.substring(0, 7)}</strong>:`,
          callbackConfirm: () => {
            let data = messageDialog.addTagData();
            AddTag(e.dataset.hash, data.name, data.type === 'annotated', data.message, data.push).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          },
        });
      },
    },
    {
      sep: true,
    },
    {
      text: "Copy Commit Hash to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.hash);
      },
    },
    {
      text: "Copy Commit Subject to Clipboard",
      callback: (e) => {
        ClipboardSetText((e.getElementsByClassName('commit__td--subject')[0] as HTMLElement).innerText);
      },
    },
  ]);

  return m;
}
