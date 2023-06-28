import { menuLabelBranch } from 'context-menus/label-branch';
import { menuLabelTag } from 'context-menus/label-tag';
import { menuLabelStash } from 'context-menus/label-stash';
import { menuLabelHead } from 'context-menus/label-head';
import { menuLabelRemoteHead } from 'context-menus/label-remote-head';
import { menuLabelRemoteBranch } from 'context-menus/label-remote-branch';
import { menuCommitRow } from 'context-menus/commit-row';
import { menuRemote } from 'context-menus/remote';
import { menuChangesFile } from 'context-menus/changes-file';
import { menuChangesList } from 'context-menus/changes-list';
import { menuBranchesList } from 'context-menus/branches-list';
import { menuBranch } from 'context-menus/branch-in-list';
import { menuRepo } from 'context-menus/repo';
import { writable } from 'svelte/store';
import { menuHash } from './hash';

interface Menus { [name: string]: Menu }

export type Menu = (e: HTMLElement) => MenuItem[];

export interface MenuItem {
  text?: string;
  // e will be the element or parent element clicked on with the menu class.
  callback?: (e: HTMLElement) => any;
  sep?: boolean;
}

const menus: Menus = {
  branch: menuLabelBranch,
  tag: menuLabelTag,
  stash: menuLabelStash,
  head: menuLabelHead,
  remoteHead: menuLabelRemoteHead,
  remoteBranch: menuLabelRemoteBranch,
  commit: menuCommitRow,
  remote: menuRemote,
  change: menuChangesFile,
  changes: menuChangesList,
  branchList: menuBranchesList,
  branchInList: menuBranch,
  repo: menuRepo,
  hash: menuHash,
};

function createContextMenu() {
  const { subscribe, set } = writable({} as MenuItem[]);

  return {
    subscribe,
    set,
    isMenu: (menu: string) => !!menus[menu],
    setMenu: (menu: string, clickedElement: HTMLElement) => set(menus[menu](clickedElement) ?? [] as MenuItem[]),
  };
}
export const contextMenu = createContextMenu();
