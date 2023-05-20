import { type Menu, type MenuItem } from "context-menus/_all";
import { parseResponse } from "scripts/parse-response";
import { currentBranch } from "stores/branches";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { PushBranch, ResetBranchToRemote, DeleteBranch, RenameBranch, RebaseOnBranch } from "wailsjs/go/main/App";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

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
        messageDialog.fillBlank({
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
    m = m.concat([
      {
        sep: true,
      },
      {
        text: "Rebase on Branch",
        callback: () => {
          RebaseOnBranch(e.dataset.branch).then(r => {
            parseResponse(r, () => {
              commitData.refresh();
              commitSignData.refresh();
            });
          });
        },
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
}
