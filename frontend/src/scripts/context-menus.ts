import { ClipboardSetText } from 'wailsjs/runtime/runtime';

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
    let m: MenuItem[] = [];
    if (e.dataset.current !== "true") {
      m.push({
        text: "Checkout Branch",
        callback: () => alert("todo: checkout"),
      });
    }
    m = m.concat([
      {
        text: "Push Branch",
        callback: () => alert("todo: push"),
      },
      {
        text: "Rename Branch",
        callback: () => alert("todo: rename"),
      },
    ]);
    if (e.dataset.current !== "true") {
      m.push({
        text: "Delete Branch",
        callback: () => alert("todo: delete"),
      });
    }
    if (e.dataset.current !== "true") {
      m = m.concat([
        {
          sep: true,
        },
        {
          text: "Rebase on Branch",
          callback: () => alert("todo: rebase"),
        },
      ]);
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

  head: (e: HTMLElement) => {
    return [
      {
        text: "ToDo",
        callback: () => alert("todo"),
      },
    ];
  },

  remoteHead: (e: HTMLElement) => {
    return [
      {
        text: "ToDo",
        callback: () => alert("todo"),
      },
      {
        sep: true,
      },
      {
        text: "Copy Remote Name to Clipboard",
        callback: (e) => {
          ClipboardSetText(e.dataset.remote);
        },
      },
    ];
  },

  remoteBranch: (e: HTMLElement) => {
    return [
      {
        text: "Push Branch",
        callback: () => alert("push"),
      },
      {
        text: "Pull Branch",
        callback: () => alert("pull"),
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
          text: "Cherry Pick Commit",
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
