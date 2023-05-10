import { ClipboardSetText } from "../../wailsjs/runtime/runtime";

interface Menus { [name: string]: Menu }

export type Menu = (e: HTMLElement) => MenuItem[];

export interface MenuItem {
  text?: string;
  // e will be the element or parent element clicked on with the menu class.
  callback?: (e: HTMLElement) => any;
  sep?: boolean;
}

export const menus: Menus = {

  branch: (e: HTMLElement) => {
    let m: MenuItem[] = [
      {
        text: "Checkout Branch",
        callback: () => alert("todo: checkout"),
      },
      {
        text: "Push Branch",
        callback: () => alert("todo: push"),
      },
      {
        text: "Rename Branch",
        callback: () => alert("todo: rename"),
      },
    ];
    if (e.dataset.active) {
      m.push({
        text: "Delete Branch",
        callback: () => alert("todo: delete"),
      });
    }
    m = m.concat([
      {
        sep: true,
      },
      {
        text: "Rebase on Branch",
        callback: () => alert("todo: rebase"),
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
    ]);

    return m;
  },

  tag: (e: HTMLElement) => {
    return [
      {
        text: "Push Tag",
        callback: () => alert("todo: push"),
      },
      {
        text: "Delete Tag",
        callback: () => alert("todo: del"),
      },
      {
        sep: true,
      },
      {
        text: "Copy Tag Name to Clipboard",
        callback: (e) => {
          ClipboardSetText(e.dataset.name);
        },
      },
    ];
  },

  stash: (e: HTMLElement) => {
    return [
      {
        text: "Apply Stash",
        callback: () => alert("todo: apply"),
      },
      {
        text: "Pop Stash",
        callback: () => alert("todo: pop"),
      },
      {
        text: "Drop Stash",
        callback: () => alert("todo: drop"),
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
        text: "Copy Stash Name to Clipboard",
        callback: (e) => {
          ClipboardSetText(window.atob(e.dataset.subject));
        },
      },
    ];
  },

  commit: (e: HTMLElement) => {
    let m: MenuItem[] = [];

    if (e.dataset.head !== 'true') {
      m = m.concat([
        {
          text: "Checkout Commit",
          callback: () => alert("todo: checkout"),
        },
        {
          text: "Cherry Pick",
          callback: () => alert("todo: checkout"),
        },
      ]);
    }

    m = m.concat([
      {
        text: "Revert Commit",
        callback: () => alert("todo: revert"),
      },
      {
        text: "Add Tag",
        callback: () => alert("todo: add tag"),
      },
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
};
