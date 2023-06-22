import switchBranch from "actions/branch/switch";
import type { Menu, MenuItem } from "context-menus/_all";
import { branches, currentBranch } from "stores/branches";
import { currentRepo, repos } from "stores/repos";
import { get } from "svelte/store";

export const menuBranchesList: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (get(repos)[get(currentRepo)].Main && get(currentBranch).Name !== get(repos)[get(currentRepo)].Main) {
    m.push({
      text: "Switch to " + get(repos)[get(currentRepo)].Main,
      callback: () => switchBranch(get(repos)[get(currentRepo)].Main),
    });
  }

  m.push({
    text: "Refresh Branches",
    callback: branches.refresh,
  });

  return m;
}
