import { type Menu, type MenuItem } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";
import switchBranch from "actions/branch/switch";
import deleteRemoteBranch from "actions/branch/remote-delete";

export const menuRemoteBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [
    {
      text: "Checkout Branch",
      callback: () => switchBranch(e.dataset.name, e.dataset.remote),
    },
    {
      text: "Delete Branch",
      callback: () => deleteRemoteBranch(e.dataset.name, e.dataset.remote),
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
  ];

  return m;
}
