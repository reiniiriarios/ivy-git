import { checkRef } from "scripts/check-ref";
import { parseResponse } from "scripts/parse-response";
import { currentBranch, type Branch } from "stores/branches";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { branchSelect } from "stores/ui";
import { CreateBranch } from "wailsjs/go/main/App";

function parseBranchPaste(input: string): string {
  // JIRA adds this to its 'create branch' feature.
  return input?.replace('git checkout -b ', '') ?? '';
}

function createBranch(hash: string = "") {
  messageDialog.confirm({
    heading: 'Create Branch',
    message: hash ? `Create a branch at commit <strong>${hash.substring(0, 7)}</strong>:` : 'Create a branch?',
    blank: "Name of Branch",
    validateBlank: checkRef,
    autoEditBlank: parseBranchPaste,
    confirm: 'Create',
    checkboxes: [{
      id: 'checkout',
      label: 'Checkout Branch',
      checked: true,
    }],
    callbackConfirm: () => {
      let name = messageDialog.blankValue();
      CreateBranch(name, hash, messageDialog.tickboxTicked('checkout')).then(r => {
        parseResponse(r, () => {
          currentBranch.set({Name: name} as Branch);
          branchSelect.set(false);
          commitData.refresh();
          commitSignData.refresh();
        })
      });
    }
  });
}

export default createBranch;
