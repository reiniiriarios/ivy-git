import { ClipboardSetText } from '../../wailsjs/runtime/runtime';

export interface Menu {
  class: string;
  items: MenuItem[];
}

export interface MenuItem {
  text?: string;
  // e will be the element or parent element clicked on with the menu class.
  callback?: (e: HTMLElement) => any;
  sep?: boolean;
}

export const menus: Menu[] = [
  {
    class: "refs__label--branch",
    items: [
      {
        text: "Checkout Branch",
        callback: () => alert('todo: checkout'),
      },
      {
        text: "Push Branch",
        callback: () => alert('todo: push'),
      },
      {
        text: "Rename Branch",
        callback: () => alert('todo: rename'),
      },
      {
        sep: true,
      },
      {
        text: "Copy Branch Name to Clipboard",
        callback: (e) => {
          ClipboardSetText(e.dataset.name)
        },
      },
    ],
  },
  {
    class: "refs__label--tag",
    items: [
      {
        text: "Checkout Commit",
        callback: () => alert('todo: checkout'),
      },
      {
        sep: true,
      },
      {
        text: "Copy Tag Name to Clipboard",
        callback: (e) => {
          ClipboardSetText(e.dataset.name)
        },
      },
    ]
  }
];
