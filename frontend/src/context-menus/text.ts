import type { Menu, MenuItem } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuText: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  let selection = window.getSelection().toString();
  m.push({
    text: "Copy",
    callback: () => {
      ClipboardSetText(selection);
    },
  });

  return m;
}
