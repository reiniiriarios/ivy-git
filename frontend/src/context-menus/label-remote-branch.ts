import { ClipboardSetText } from "wailsjs/runtime/runtime";

import type { Menu, MenuItem } from "context-menus/_all";

import switchBranch from "actions/branch/switch";
import pullBranch from "actions/branch/pull";
import deleteRemoteBranch from "actions/branch/remote-delete";

export const menuLabelRemoteBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => switchBranch(e.dataset.branch, e.dataset.remote),
    });
  }

  m = m.concat([
    {
      text: "Pull Branch",
      callback: () => pullBranch(e.dataset.branch, e.dataset.remote),
    },
    {
      text: "Delete Remote Branch",
      callback: () => deleteRemoteBranch(e.dataset.branch, e.dataset.remote),
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
  ]);

  return m;
}
