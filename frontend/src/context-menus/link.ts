import type { Menu, MenuItem } from "context-menus/_all";
import { BrowserOpenURL, ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuLink: Menu = (e: HTMLAnchorElement) => {
  let m: MenuItem[] = [
    {
      text: "Open Link",
      callback: () => {
        console.log("Opening link:", e.innerText);
        BrowserOpenURL(e.href);
      },
    },
    {
      text: "Copy URL to Clipboard",
      callback: () => {
        ClipboardSetText(e.href);
      },
    }
  ];

  if (e.href.startsWith('mailto:')) {
    m.push({
      text: "Copy Email Address to Clipboard",
      callback: () => {
        ClipboardSetText(e.href.substring(7).trim());
      },
    })
  }

  if (e.href.startsWith('tel:')) {
    m.push({
      text: "Copy Phone Number to Clipboard",
      callback: () => {
        ClipboardSetText(e.href.substring(4).trim());
      },
    })
  }

  return m;
}
