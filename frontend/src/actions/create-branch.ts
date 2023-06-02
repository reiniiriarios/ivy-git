import { checkRef } from "scripts/check-ref";
import { parseResponse } from "scripts/parse-response";
import { currentBranch, type Branch } from "stores/branches";
import { commitData, commitSignData } from "stores/commit-data";
import { messageDialog } from "stores/message-dialog";
import { CreateBranch } from "wailsjs/go/main/App";

function createBranch(hash: string) {
  messageDialog.confirm({
    heading: 'Create Branch',
    message: `Create a branch at commit <strong>${hash.substring(0, 7)}</strong>:`,
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
        hash,
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

export default createBranch;
