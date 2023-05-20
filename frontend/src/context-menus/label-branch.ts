import { type Menu, type MenuItem } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { currentBranch } from "stores/branches";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { PushBranch, ResetBranchToRemote } from "wailsjs/go/main/App";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuLabelBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];
  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => {
        currentBranch.set(e.dataset.name);
      },
    });
  }
  m = m.concat([
    {
      text: "Push Branch",
      callback: (e) => {
        PushBranch(e.dataset.name).then(r => {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        })
      },
    },
    {
      text: "Rename Branch",
      callback: () => alert("todo: rename"),
    },
  ]);
  if (e.dataset.current !== "true") {
    m.push({
      text: "Delete Branch",
      callback: () => alert("todo: delete"),
    });
  }

  if (e.dataset.upstream) {
    m.push({
      text: "Reset Local Branch to Remote",
      callback: (e) => {
        messageDialog.confirm({
          heading: 'Reset Local Branch to Remote',
          message: `Are you sure you want to reset the local branch <strong>${e.dataset.branch}</strong> to its default remote?`,
          confirm: 'Reset',
          okay: 'Cancel',
          callbackConfirm: () => {
            ResetBranchToRemote(e.dataset.branch).then(() => {
              commitData.refresh();
              commitSignData.refresh();
            });
          },
        });
      },
    });
  }

  if (e.dataset.current !== "true") {
    m = m.concat([
      {
        sep: true,
      },
      {
        text: "Rebase on Branch",
        callback: () => alert("todo: rebase"),
      },
    ]);
  }
  m = m.concat([
    {
      sep: true,
    },
    {
      text: "Copy Branch Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
  ]);

  return m;
}
