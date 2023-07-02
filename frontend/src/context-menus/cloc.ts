import type { Menu, MenuItem } from "context-menus/_all";
import { cloc } from "stores/cloc";
import { messageDialog } from "stores/message-dialog";

export const menuCloc: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  m.push({
    text: "Clear Code Breakdown Data",
    callback: (e) => {
      messageDialog.confirm({
        heading: 'Clear Code Breakdown Data',
        message: 'Are you sure you want to delete code breakdown data for this repo?',
        confirm: 'Clear',
        callbackConfirm: () => cloc.reset(),
      });
    },
  });

  return m;
}
