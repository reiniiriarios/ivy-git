import { type Menu, type MenuItem } from "context-menus/_all";

import { ClipboardSetText } from "wailsjs/runtime/runtime";
import {
  PushBranch,
  ResetBranchToRemote,
  DeleteBranch,
  RenameBranch,
  RebaseOnBranch,
  MergeCommit,
  MergeRebase,
  MergeSquash,
  MergeFastForward
} from "wailsjs/go/main/App";

import { get } from 'svelte/store';
import { currentBranch } from "stores/branches";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { settings } from "stores/settings";

import { parseResponse } from "scripts/parse-response";

export const menuLabelBranch: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];
  if (e.dataset.current !== "true") {
    m.push({
      text: "Checkout Branch",
      callback: () => {
        currentBranch.set(e.dataset.name);
      },
    });
  }
  m = m.concat([
    {
      text: "Push Branch",
      callback: (e) => {
        PushBranch(e.dataset.name).then(r => {
          parseResponse(r, () => {
            commitData.refresh();
            commitSignData.refresh();
          });
        })
      },
    },
    {
      text: "Rename Branch",
      callback: () => {
        messageDialog.confirm({
          heading: 'Rename Branch',
          message: `Rename <strong>${e.dataset.branch}</strong> locally and on all remotes to:`,
          confirm: 'Rename',
          blank: 'New Name',
          okay: 'Cancel',
          callbackConfirm: () => {
            RenameBranch(e.dataset.branch, messageDialog.blankValue()).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          },
        });
      },
    },
  ]);
  if (e.dataset.current !== "true") {
    m.push({
      text: "Delete Branch",
      callback: () => {
        let opts = [{id: 'force', label: 'Force Delete'}];
        if (e.dataset.upstream) {
          opts.push({id: 'remote', label: 'Delete on Remote'});
        }
        messageDialog.confirm({
          heading: 'Delete Branch',
          message: `Are you sure you want to delete the branch <strong>${e.dataset.branch}</strong>?`,
          confirm: 'Delete',
          okay: 'Cancel',
          checkboxes: opts,
          callbackConfirm: () => {
            DeleteBranch(
              e.dataset.branch,
              messageDialog.tickboxTicked('force'),
              messageDialog.tickboxTicked('remote')
            ).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          },
        });
      },
    });
  }

  if (e.dataset.upstream) {
    m.push({
      text: "Reset Local Branch to Remote",
      callback: (e) => {
        messageDialog.confirm({
          heading: 'Reset Local Branch to Remote',
          message: `Are you sure you want to reset the local branch <strong>${e.dataset.branch}</strong> to its default remote?`,
          confirm: 'Reset',
          okay: 'Cancel',
          callbackConfirm: () => {
            ResetBranchToRemote(e.dataset.branch).then(() => {
              commitData.refresh();
              commitSignData.refresh();
            });
          },
        });
      },
    });
  }

  if (e.dataset.current !== "true") {
    m.push({
      sep: true,
    });

    let workflow = get(settings).Workflow;
    if (workflow === 'rebase' || workflow === 'squash') {
      m.push({
        text: "Rebase on Branch",
        callback: () => {
          RebaseOnBranch(e.dataset.branch).then(r => {
            parseResponse(r, () => {
              commitData.refresh();
              commitSignData.refresh();
            });
          });
        },
      });
    }
    if (workflow === 'rebase') {
      m.push({
        text: "Rebase and Merge into Current Branch",
        callback: () => {
          messageDialog.confirm({
            heading: 'Rebase and Merge into Current Branch',
            message: `Rebase <strong>${e.dataset.branch}</strong> onto current branch and merge?`,
            confirm: 'Merge',
            okay: 'Cancel',
            callbackConfirm: () => {
              MergeRebase(e.dataset.branch).then(r => {
                parseResponse(r, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                });
              });
            },
          });
        },
      });
      m.push({
        text: "Fast-forward Merge",
        callback: () => {
          messageDialog.confirm({
            heading: 'Fast-forward Merge',
            message: `Merge the current branch into <strong>${e.dataset.branch}</strong> via fast-forward only?`,
            confirm: 'Merge',
            okay: 'Cancel',
            callbackConfirm: () => {
              MergeFastForward(e.dataset.branch).then(r => {
                parseResponse(r, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                });
              });
            },
          });
        },
      });
    }
    if (workflow === 'squash' || workflow === 'merge') {
      m.push({
        text: "Merge into Current Branch",
        callback: () => {
          messageDialog.confirm({
            heading: 'Merge into Current Branch',
            message: `Merge <strong>${e.dataset.branch}</strong> into current branch?`,
            confirm: 'Merge',
            okay: 'Cancel',
            checkboxes: [
              {
                id: 'no-ff',
                label: 'Create a new commit even if fast-forward is possible',
                checked: true,
              },
              {
                id: 'no-commit',
                label: 'No Commit',
                checked: false,
              },
            ],
            callbackConfirm: () => {
              MergeCommit(
                e.dataset.branch,
                messageDialog.tickboxTicked('no-commit'),
                messageDialog.tickboxTicked('no-ff')
              ).then(r => {
                parseResponse(r, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                });
              });
            },
          });
        },
      });
    }
    if (workflow === 'squash') {
      m.push({
        text: "Squash & Merge onto Current Branch",
        callback: () => {
          messageDialog.confirm({
            heading: 'Squash & Merge onto Current Branch',
            message: `Squash <strong>${e.dataset.branch}</strong> and merge onto current branch?`,
            confirm: 'Merge',
            okay: 'Cancel',
            callbackConfirm: () => {
              MergeSquash(e.dataset.branch).then(r => {
                parseResponse(r, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                });
              });
            },
          });
        },
      });
    }
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
}
