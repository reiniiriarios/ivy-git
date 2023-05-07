export interface Menu {
  class: string;
  items: MenuItem[];
}

export interface MenuItem {
  text?: string;
  callback?: () => any;
  sep?: boolean;
}

export const menus: Menu[] = [
  {
    class: "refs__branch",
    items: [
      {
        text: "Checkout Branch",
        callback: () => alert('oh hey'),
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
        callback: () => alert('todo: copy'),
      },
    ],
  },
];
