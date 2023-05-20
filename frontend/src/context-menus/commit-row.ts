import type { Menu, MenuItem } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuCommitRow: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.head !== 'true') {
    m = m.concat([
      {
        text: "Checkout Commit",
        callback: () => alert("todo: checkout"),
      },
      {
        text: "Cherry Pick Commit",
        callback: () => alert("todo: checkout"),
      },
    ]);
  }

  m = m.concat([
    {
      text: "Revert Commit",
      callback: () => alert("todo: revert"),
    },
    {
      text: "Add Tag",
      callback: () => alert("todo: add tag"),
    },
    {
      sep: true,
    },
    {
      text: "Copy Commit Hash to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.hash);
      },
    },
    {
      text: "Copy Commit Subject to Clipboard",
      callback: (e) => {
        ClipboardSetText((e.getElementsByClassName('commit__td--subject')[0] as HTMLElement).innerText);
      },
    },
  ]);

  return m;
}
