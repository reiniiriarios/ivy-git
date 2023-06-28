import type { Menu, MenuItem } from "context-menus/_all";
import { ClipboardGetText, ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuInput: Menu = (e: HTMLInputElement | HTMLTextAreaElement) => {
  let m: MenuItem[] = [];

  let selection = window.getSelection().toString();
  m.push({
    text: "Copy",
    callback: () => {
      ClipboardSetText(selection);
    },
  });

  m.push({
    text: "Cut",
    callback: () => {
      ClipboardSetText(selection);
      e.value = '';
    },
  });

  m.push({
    text: "Paste",
    callback: () => {
      ClipboardGetText().then(text => e.value = text);
    },
  });

  return m;
}
