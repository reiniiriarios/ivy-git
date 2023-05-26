import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { CreateBranch, AddTag, CheckoutCommit, RevertCommit, HardReset, SoftReset, CherryPick } from "wailsjs/go/main/App";

import { get } from 'svelte/store';

import type { Menu, MenuItem } from "context-menus/_all";

import { parseResponse } from "scripts/parse-response";

import { commitData, commitSignData, HEAD } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { checkRef } from "scripts/check-ref";
import { currentBranch, type Branch, detachedHead } from "stores/branches";
import { inProgressCommitMessage } from "stores/ui";

export const menuCommitRow: Menu = (e: HTMLElement) => {
  let m: MenuItem[] = [];

  m = m.concat([
    {
      text: "Create Branch",
      callback: () => {
        messageDialog.confirm({
          heading: 'Create Branch',
          message: `Create a branch at commit <strong>${e.dataset.hash.substring(0, 7)}</strong>:`,
          blank: "Name of Branch",
          validateBlank: checkRef,
          confirm: 'Create',
          checkboxes: [{
            id: 'checkout',
            label: 'Checkout Branch',
            checked: true,
          }],
          callbackConfirm: () => {
            CreateBranch(
              messageDialog.blankValue(),
              e.dataset.hash,
              messageDialog.tickboxTicked('checkout')
            ).then(r => {
              parseResponse(r, () => {
                currentBranch.set({Name: messageDialog.blankValue()} as Branch);
                commitData.refresh();
                commitSignData.refresh();
              })
            });
          }
        });
      }
    },
    {
      text: "Add Tag",
      callback: () => {
        messageDialog.addTag({
          message: `Add tag to commit <strong>${e.dataset.hash.substring(0, 7)}</strong>:`,
          callbackConfirm: () => {
            let data = messageDialog.addTagData();
            AddTag(e.dataset.hash, data.name, data.type === 'annotated', data.message, data.push).then(r => {
              parseResponse(r, () => {
                commitData.refresh();
                commitSignData.refresh();
              });
            });
          },
        });
      },
    },
    {
      sep: true,
    }
  ]);

  if (e.dataset.head !== 'true') {
    m = m.concat([
      {
        text: "Checkout Commit",
        callback: () => {
          let co = () => CheckoutCommit(e.dataset.hash).then(result => {
            parseResponse(result, () => {
              commitData.refresh();
              commitSignData.refresh();
              currentBranch.detach();
            });
          });
          if (get(detachedHead)) {
            messageDialog.confirm({
              heading: 'Checkout Commit',
              message: 'You are currently in a <strong>detached HEAD</strong> state. Checking out a different commit could result in lost work. Continue?',
              confirm: 'Checkout',
              callbackConfirm: co,
            });
          }
          else {
            co();
          }
        },
      },
    ]);
  }

  if (e.dataset.merge !== 'true') {
    m = m.concat([
      {
        text: "Cherry Pick Commit",
        callback: () => {
          messageDialog.confirm({
            heading: 'Cherry Pick Commit',
            message: `Cherry pick commit <strong>${e.dataset.hash.substring(0, 7)}</strong>.`,
            checkboxes: [
              {
                id: 'record',
                label: 'Record Original Hash',
              },
              {
                id: 'no_commit',
                label: 'No Commit',
              },
            ],
            confirm: 'Cherry Pick',
            callbackConfirm: () => {
              let no_commit = messageDialog.tickboxTicked('no_commit');
              CherryPick(e.dataset.hash, messageDialog.tickboxTicked('record'), no_commit).then(result => {
                parseResponse(result, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                  if (no_commit) {
                    inProgressCommitMessage.fetch();
                  }
                }, () => {
                  inProgressCommitMessage.fetch();
                });
              });
            }
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
      text: "Revert Commit",
      callback: () => {
        messageDialog.confirm({
          heading: 'Revert Commit',
          message: `Are you sure you want to revert <strong>${e.dataset.hash.substring(0, 7)}</strong>?`,
          confirm: 'Revert',
          callbackConfirm: () => {
            RevertCommit(e.dataset.hash).then(result => {
              parseResponse(result, () => {
                commitData.refresh();
                commitSignData.refresh();
                currentBranch.clear();
              });
            });
          },
        });
      },
    },
    // todo: do not enable until there's a safety check dumdum
    // {
    //   text: "Drop Commit",
    //   callback: () => {
    //     messageDialog.confirm({
    //       heading: 'Drop Commit',
    //       message: `Are you sure you want to drop <strong>${e.dataset.hash.substring(0, 7)}</strong>?`,
    //       confirm: 'Drop',
    //       callbackConfirm: () => {
    //         DropCommit(e.dataset.hash).then(result => {
    //           parseResponse(result, () => {
    //             commitData.refresh();
    //             commitSignData.refresh();
    //             currentBranch.clear();
    //           });
    //         });
    //       },
    //     });
    //   },
    // }
  ]);

  if (e.dataset.hash !== get(HEAD).Hash) {
    m = m.concat([
      {
        text: 'Soft Reset to This Commit',
        callback: () => {
          messageDialog.confirm({
            heading: 'Soft Reset',
            message: `Are you sure you want to soft reset to <strong>${e.dataset.hash.substring(0, 7)}</strong>?`,
            confirm: 'Soft Reset',
            callbackConfirm: () => {
              SoftReset(e.dataset.hash).then(result => {
                parseResponse(result, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                  currentBranch.clear();
                });
              });
            },
          });
        },
      },
      {
        text: 'Hard Reset to This Commit',
        callback: () => {
          messageDialog.confirm({
            heading: 'Hard Reset',
            message: `Are you sure you want to hard reset to <strong>${e.dataset.hash.substring(0, 7)}</strong>?`,
            confirm: 'Hard Reset',
            callbackConfirm: () => {
              HardReset(e.dataset.hash).then(result => {
                parseResponse(result, () => {
                  commitData.refresh();
                  commitSignData.refresh();
                  currentBranch.clear();
                });
              });
            },
          });
        },
      }
    ]);
  }

  m = m.concat([
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
