import type { Menu, MenuItem } from "context-menus/_all";
import { currentCommit } from "stores/commit-details";
import { commits, commitsMap } from "stores/commits";
import { currentTab } from "stores/ui";
import { get } from "svelte/store";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuHash: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.linked) {
    m.push({
      text: "Go to Commit",
      callback: () => {
        currentTab.set('tree');
        currentCommit.set(get(commits)[get(commitsMap).get(e.dataset.hash)]);
      },
    });
  }

  m.push({
    text: "Copy Hash to Clipboard",
    callback: (e) => {
      ClipboardSetText(e.dataset.hash);
    },
  });

  return m;
}
