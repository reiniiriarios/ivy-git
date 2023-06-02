import { ClipboardSetText } from "wailsjs/runtime/runtime";

import type { Menu, MenuItem } from "context-menus/_all";

import checkoutBranch from "actions/checkout-branch";
import pullBranch from "actions/pull-branch";
import deleteRemoteBranch from "actions/delete-remote-branch";

export const menuLabelRemoteBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => checkoutBranch(e.dataset.branch, e.dataset.remote),
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
