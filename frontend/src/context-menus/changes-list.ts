import { stageAll } from "actions/stage/stage-all";
import { unstageAll } from "actions/stage/unstage-all";
import type { Menu, MenuItem } from "context-menus/_all";
import { changes } from "stores/changes";

export const menuChangesList: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.staged === "true") {
    m.push({
      text: "Unstage All",
      callback: unstageAll,
    });
  }

  if (e.dataset.unstaged === "true") {
    m.push({
      text: "Stage All",
      callback: stageAll,
    });
  }

  m.push({
    text: "Refresh Changes",
    callback: changes.refresh,
  });

  return m;
}
