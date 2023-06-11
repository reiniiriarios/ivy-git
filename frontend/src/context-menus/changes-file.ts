import { stageFile } from "actions/stage/stage";
import { unstageFile } from "actions/stage/unstage";
import type { Menu, MenuItem } from "context-menus/_all";
import { changes, type Change } from "stores/changes";
import { get } from "svelte/store";
import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { discardChanges } from "actions/stage/discard-changes";

export const menuChangesFile: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.staged === "true") {
    m = m.concat([
      {
        text: "Unstage File",
        callback: () => {
          let file: Change = get(changes).x[e.dataset.file];
          if (file) {
            unstageFile(file, e.dataset.partial === 'true');
          }
        },
      },
      {
        sep: true,
      },
    ]);
  }
  else if (e.dataset.conflict !== "true") {
    m = m.concat([
      {
        text: "Stage File",
        callback: () => {
          let file: Change = get(changes).y[e.dataset.file];
          if (file) {
            stageFile(file, e.dataset.partial === 'true');
          }
        },
      },
      {
        sep: true,
      },
      {
        text: "Discard Changes",
        callback: () => discardChanges(e.dataset.file),
      },
      {
        sep: true,
      },
    ]);
  }

  m = m.concat([
    {
      text: "Copy Filename to Clipboard",
      callback: () => {
        let f = e.dataset.file.replaceAll('\\','/').substring(e.dataset.file.lastIndexOf('/') + 1);
        ClipboardSetText(f);
      },
    },
    {
      text: "Copy File Path to Clipboard",
      callback: () => ClipboardSetText(e.dataset.file),
    }
  ]);

  return m;
}
