import type { Menu, MenuItem } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { PullRemoteBranch, DeleteRemoteBranch } from "wailsjs/go/main/App";
import { ClipboardSetText } from "wailsjs/runtime/runtime";


export const menuLabelRemoteBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [
    {
      text: "Pull Branch",
      callback: () => {
        PullRemoteBranch(e.dataset.remote, e.dataset.branch, true).then(r => {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        })
      },
    },
    {
      text: "Delete Remote Branch",
      callback: () => {
        let opts = [{id: 'force', label: 'Force Delete'}];
        messageDialog.confirm({
          heading: 'Delete Remote Branch',
          message: `Are you sure you want to delete the remote branch <strong>${e.dataset.remote}/${e.dataset.branch}</strong>?`,
          confirm: 'Delete',
          okay: 'Cancel',
          checkboxes: opts,
          callbackConfirm: () => {
            DeleteRemoteBranch(
              e.dataset.branch,
              e.dataset.remote,
              messageDialog.tickboxTicked('force')
            ).then(r => {
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
      text: "Copy Branch Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
    {
      text: "Copy Remote Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.remote);
      },
    },
  ];

  return m;
}
