import type { Menu } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuLabelRemoteHead: Menu = (e: HTMLElement) => {
  return [
    {
      text: "ToDo",
      callback: () => alert("todo"),
    },
    {
      sep: true,
    },
    {
      text: "Copy Remote Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.remote);
      },
    },
  ];
}
