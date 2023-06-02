import { currentBranch, detachedHead } from "stores/branches";
import { messageDialog } from "stores/message-dialog";
import { get } from "svelte/store";

function checkoutBranch(branch: string, remote: string = "") {
  if (get(detachedHead)) {
    messageDialog.confirm({
      heading: 'Checkout Branch',
      message: 'You are currently in a <strong>detached HEAD</strong> state. Checking out a branch could result in lost work. Continue?',
      confirm: 'Checkout',
      callbackConfirm: () => currentBranch.switch(branch, remote),
    });
  }
  else {
    currentBranch.switch(branch, remote);
  }
}

export default checkoutBranch;
