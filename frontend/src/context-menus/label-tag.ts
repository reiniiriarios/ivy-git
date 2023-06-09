import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { type Menu } from "context-menus/_all";
import pushTag from "actions/tag/push";
import deleteTag from "actions/tag/delete";

export const menuLabelTag: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Push Tag",
      callback: () => pushTag(e.dataset.name),
    },
    {
      text: "Delete Tag",
      callback: () => deleteTag(e.dataset.name),
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
