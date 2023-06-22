import { type Menu, type MenuItem } from "context-menus/_all";
import { get } from "svelte/store";
import { finderWord } from "stores/env";
import { currentRepo } from "stores/repos";
import { OpenRepoInFinder } from "wailsjs/go/main/App";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export const menuRepo: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  m = m.concat([
    {
      text: "Refresh Repo",
      callback: currentRepo.refresh,
    },
    {
      text: `Open Repo Directory in ${get(finderWord)}`,
      callback: (e) => {
        OpenRepoInFinder(e.dataset.id);
      },
    },
    {
      sep: true,
    },
    {
      text: "Copy Repo Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
    {
      text: "Copy Repo Path to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.path);
      },
    },
  ]);

  return m;
}
