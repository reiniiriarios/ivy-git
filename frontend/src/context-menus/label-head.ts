import type { Menu } from "context-menus/_all";

export const menuLabelHead: Menu = (e: HTMLElement) => {
  return [
    {
      text: "ToDo",
      callback: () => alert("todo"),
    },
  ];
}
