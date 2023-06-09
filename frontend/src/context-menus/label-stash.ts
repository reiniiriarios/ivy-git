import applyStash from "actions/stash/apply";
import createBranchFromStash from "actions/stash/create-branch";
import dropStash from "actions/stash/drop";
import popStash from "actions/stash/pop";
import type { Menu } from "context-menus/_all";
import { ClipboardSetText } from "wailsjs/runtime/runtime";


export const menuLabelStash: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Apply Stash",
      callback: (e) => applyStash(e.dataset.ref),
    },
    {
      text: "Pop Stash",
      callback: (e) => popStash(e.dataset.ref, e.dataset.hash),
    },
    {
      text: "Drop Stash",
      callback: (e) => dropStash(e.dataset.ref, e.dataset.hash),
    },
    {
      sep: true,
    },
    {
      text: "Create Branch from Stash",
      callback: (e) => createBranchFromStash(e.dataset.ref, e.dataset.hash),
    },
    {
      sep: true,
    },
    {
      text: "Copy Stash Hash to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.hash);
      },
    },
    {
      text: "Copy Stash Ref to Clipboard",
      callback: (e) => {
        ClipboardSetText(e.dataset.ref);
      },
    },
    {
      text: "Copy Stash Name to Clipboard",
      callback: (e) => {
        ClipboardSetText(window.atob(e.dataset.subject));
      },
    },
  ];
}
