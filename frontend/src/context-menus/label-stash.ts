import type { Menu } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";


export const menuLabelStash: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Apply Stash",
      callback: () => alert("todo: apply"),
    },
    {
      text: "Pop Stash",
      callback: () => alert("todo: pop"),
    },
    {
      text: "Drop Stash",
      callback: () => alert("todo: drop"),
    },
    {
      sep: true,
    },
    {
      text: "Copy Stash Hash to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.hash);
      },
    },
    {
      text: "Copy Stash Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(window.atob(e.dataset.subject));
      },
    },
  ];
}
