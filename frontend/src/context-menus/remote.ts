import deleteRemote from "actions/remote/delete";
import type { Menu } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuRemote: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Delete Remote",
      callback: () => deleteRemote(e.dataset.name),
    },
    {
      text: "Copy Remote URL",
      callback: () => ClipboardSetText(e.dataset.url),
    },
  ];
}
