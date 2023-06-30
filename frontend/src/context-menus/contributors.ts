import type { Menu, MenuItem } from "context-menus/_all";
import { contributors } from "stores/contributors";
import { messageDialog } from "stores/message-dialog";

export const menuContributors: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  m.push({
    text: "Clear Contributor Data",
    callback: (e) => {
      messageDialog.confirm({
        heading: 'Clear Contributor Data',
        message: 'Are you sure you want to delete contributor data for this repo?',
        confirm: 'Clear',
        callbackConfirm: () => contributors.reset(),
      });
    },
  });

  return m;
}
