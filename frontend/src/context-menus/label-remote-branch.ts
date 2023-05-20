import type { Menu } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { PullRemoteBranch } from "wailsjs/go/main/App";
import { ClipboardSetText } from "wailsjs/runtime/runtime";


export const menuLabelRemoteBranch: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Pull Branch",
      callback: () => {
        PullRemoteBranch(e.dataset.remote, e.dataset.branch, true).then(r => {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        })
      },
    },
    {
      text: "Reset Local Branch to Remote",
      callback: () => alert("reset"),
    },
    {
      sep: true,
    },
    {
      text: "Copy Branch Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
    {
      text: "Copy Remote Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.remote);
      },
    },
  ];
}
