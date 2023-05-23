import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { CreateBranch, AddTag, CheckoutCommit, RevertCommit, HardReset } from "wailsjs/go/main/App";

import { get } from 'svelte/store';

import type { Menu, MenuItem } from "context-menus/_all";

import { parseResponse } from "scripts/parse-response";

import { commitData, commitSignData, HEAD } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { checkRef } from "scripts/check-ref";
import { currentBranch } from "stores/branches";

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
                commitData.refresh();
                commitSignData.refresh();
              })
            });
          }
        })
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
          CheckoutCommit(e.dataset.hash).then(result => {
            parseResponse(result, () => {
              commitData.refresh();
              commitSignData.refresh();
              currentBranch.clear();
            });
          });
        },
      },
      {
        text: "Cherry Pick Commit",
        callback: () => {},
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
    }
  ]);

  if (e.dataset.hash !== get(HEAD).Hash) {
    m.push({
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
    });
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
