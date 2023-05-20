import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { type Menu } from "context-menus/_all";

export const menuLabelTag: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Push Tag",
      callback: () => alert("todo: push"),
    },
    {
      text: "Delete Tag",
      callback: () => alert("todo: del"),
    },
    {
      sep: true,
    },
    {
      text: "Copy Tag Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
  ];
}
