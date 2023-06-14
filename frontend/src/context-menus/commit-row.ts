import { ClipboardSetText } from "wailsjs/runtime/runtime";

import { get } from 'svelte/store';

import type { Menu, MenuItem } from "context-menus/_all";

import { HEAD } from "stores/commits";
import createBranch from "actions/branch/create";
import addTag from "actions/tag/add";
import checkoutCommit from "actions/commit/checkout";
import cherryPick from "actions/commit/cherry-pick";
import revertCommit from "actions/commit/revert";
import softReset from "actions/commit/reset-soft";
import hardReset from "actions/commit/reset-hard";
import dropCommit from "actions/commit/drop";

export const menuCommitRow: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  m = m.concat([
    {
      text: "Create Branch",
      callback: () => createBranch(e.dataset.hash),
    },
    {
      text: "Add Tag",
      callback: () => addTag(e.dataset.hash),
    },
    {
      sep: true,
    }
  ]);

  if (e.dataset.head !== 'true') {
    m = m.concat([
      {
        text: "Checkout Commit",
        callback: () => checkoutCommit(e.dataset.hash),
      },
    ]);
  }

  if (e.dataset.merge !== 'true') {
    m = m.concat([
      {
        text: "Cherry Pick Commit",
        callback: () => cherryPick(e.dataset.hash),
      },
    ]);
  }

  m = m.concat([
    {
      sep: true,
    },
    {
      text: "Revert Commit",
      callback: () => revertCommit(e.dataset.hash),
    },
    {
      // todo: safety check for this one, woof.
      text: "Drop Commit",
      callback: () => dropCommit(e.dataset.hash),
    }
  ]);

  if (e.dataset.hash !== get(HEAD).Hash) {
    m = m.concat([
      {
        text: 'Soft Reset to This Commit',
        callback: () => softReset(e.dataset.hash),
      },
      {
        text: 'Hard Reset to This Commit',
        callback: () => hardReset(e.dataset.hash),
      }
    ]);
  }

  m = m.concat([
    {
      sep: true,
    },
    {
      text: "Copy Commit Hash to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.hash);
      },
    },
    {
      text: "Copy Commit Subject to Clipboard",
      callback: (e) => {
        ClipboardSetText((e.getElementsByClassName('commit__td--subject')[0] as HTMLElement).innerText);
      },
    },
  ]);

  return m;
}
