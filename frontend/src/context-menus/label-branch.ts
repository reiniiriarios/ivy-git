import { type Menu, type MenuItem } from "context-menus/_all";

import { ClipboardSetText } from "wailsjs/runtime/runtime";

import { get } from 'svelte/store';

import { detachedHead } from "stores/branches";
import { settings } from "stores/settings";

import deleteBranch from "actions/delete-branch";
import pushBranch from "actions/push-branch";
import renameBranch from "actions/rename-branch";
import checkoutBranch from "actions/checkout-branch";
import resetBranchToRemote from "actions/reset-branch-to-remote";
import rebaseOnBranch from "actions/rebase-on-branch";
import rebaseAndMergeIntoCurrentBranch from "actions/rebase-merge-into-current";
import fastForwardMerge from "actions/fast-forward-merge";
import mergeBranch from "actions/merge-branch";
import squashMergeBranch from "actions/squash-merge-branch";

export const menuLabelBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];
  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => checkoutBranch(e.dataset.name),
    });
  }
  m = m.concat([
    {
      text: "Push Branch",
      callback: (e) => pushBranch(e.dataset.name, e.dataset.branch),
    },
    {
      text: "Rename Branch",
      callback: () => renameBranch(e.dataset.branch),
    },
  ]);
  if (e.dataset.current !== "true") {
    m.push({
      text: "Delete Branch",
      callback: () => deleteBranch(e.dataset.branch, !!e.dataset.upstream),
    });
  }

  if (e.dataset.upstream) {
    m.push({
      text: "Reset Local Branch to Remote",
      callback: (e) => resetBranchToRemote(e.dataset.branch),
    });
  }

  if (e.dataset.current !== "true" && !get(detachedHead)) {
    m.push({
      sep: true,
    });

    let workflow = get(settings).Workflow;
    if (workflow === 'rebase' || workflow === 'squash') {
      m.push({
        text: "Rebase on Branch",
        callback: () => rebaseOnBranch(e.dataset.branch),
      });
    }
    if (workflow === 'rebase') {
      m.push({
        text: "Rebase and Merge into Current Branch",
        callback: () => rebaseAndMergeIntoCurrentBranch(e.dataset.branch),
      });
      m.push({
        text: "Fast-forward Merge",
        callback: () => fastForwardMerge(e.dataset.branch),
      });
    }
    if (workflow === 'squash' || workflow === 'merge') {
      m.push({
        text: "Merge into Current Branch",
        callback: () => mergeBranch(e.dataset.branch),
      });
    }
    if (workflow === 'squash') {
      m.push({
        text: "Squash & Merge onto Current Branch",
        callback: () => squashMergeBranch(e.dataset.branch),
      });
    }
  }
  m = m.concat([
    {
      sep: true,
    },
    {
      text: "Copy Branch Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.name);
      },
    },
  ]);

  return m;
}
