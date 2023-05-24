import { ClipboardSetText } from "src/_tmp";
import { PullRemoteBranch, DeleteRemoteBranch } from "src/_tmp";

import { get } from "svelte/store";

import type { Menu, MenuItem } from "context-menus/_all";

import { parseResponse } from "scripts/parse-response";

import { currentBranch } from "stores/branches";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";


export const menuLabelRemoteBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => {
        if (get(currentBranch).Name === 'HEAD') {
          messageDialog.confirm({
            heading: 'Checkout Branch',
            message: 'You are currently in a <strong>detached HEAD</strong> state. Checking out a branch could result in lost work. Continue?',
            confirm: 'Checkout',
            callbackConfirm: () => currentBranch.switch(e.dataset.branch, e.dataset.remote),
          });
        }
        else {
          currentBranch.switch(e.dataset.branch, e.dataset.remote);
        }
      },
    });
  }

  m = m.concat([
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
      text: "Delete Remote Branch",
      callback: () => {
        let opts = [{id: 'force', label: 'Force Delete'}];
        messageDialog.confirm({
          heading: 'Delete Remote Branch',
          message: `Are you sure you want to delete the remote branch <strong>${e.dataset.remote}/${e.dataset.branch}</strong>?`,
          confirm: 'Delete',
          okay: 'Cancel',
          checkboxes: opts,
          callbackConfirm: () => {
            DeleteRemoteBranch(
              e.dataset.branch,
              e.dataset.remote,
              messageDialog.tickboxTicked('force')
            ).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          },
        });
      },
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
  ]);

  return m;
}
