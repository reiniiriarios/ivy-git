import deleteRemote from "actions/remote/delete";
import type { Menu } from "context-menus/_all";

export const menuRemote: Menu = (e: HTMLElement) => {
  return [
    {
      text: "Delete Remote",
      callback: () => deleteRemote(e.dataset.name),
    },
  ];
}
