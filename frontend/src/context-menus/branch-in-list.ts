import { type Menu, type MenuItem } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";
import deleteBranch from "actions/branch/delete";
import renameBranch from "actions/branch/rename";
import checkoutBranch from "actions/branch/checkout";

export const menuBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => checkoutBranch(e.dataset.name),
    });
  }

  m = m.concat([
    {
      text: "Rename Branch",
      callback: () => renameBranch(e.dataset.name),
    },
  ]);

  if (e.dataset.current !== "true") {
    m.push({
      text: "Delete Branch",
      callback: () => deleteBranch(e.dataset.name, !!e.dataset.upstream),
    });
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
