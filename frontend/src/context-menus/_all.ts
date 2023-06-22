import { menuLabelBranch } from 'context-menus/label-branch';
import { menuLabelTag } from 'context-menus/label-tag';
import { menuLabelStash } from 'context-menus/label-stash';
import { menuLabelHead } from 'context-menus/label-head';
import { menuLabelRemoteHead } from 'context-menus/label-remote-head';
import { menuLabelRemoteBranch } from 'context-menus/label-remote-branch';
import { menuCommitRow } from 'context-menus/commit-row';
import { menuRemote } from './remote';
import { menuChangesFile } from './changes-file';
import { menuChangesList } from './changes-list';
import { menuBranchesList } from './branches-list';
import { menuBranch } from './branch-in-list';

interface Menus { [name: string]: Menu }

export type Menu = (e: HTMLElement) => MenuItem[];

export interface MenuItem {
  text?: string;
  // e will be the element or parent element clicked on with the menu class.
  callback?: (e: HTMLElement) => any;
  sep?: boolean;
}

export const menus: Menus = {
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
};
