import { parseResponse } from "scripts/parse-response";
import { currentBranch, detachedHead } from "stores/branches";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { get } from "svelte/store";
import { CheckoutCommit } from "wailsjs/go/main/App";

function checkoutCommit(hash: string) {
  if (get(detachedHead)) {
    messageDialog.confirm({
      heading: 'Checkout Commit',
      message: 'You are currently in a <strong>detached HEAD</strong> state. Checking out a different commit could result in lost work. Continue?',
      confirm: 'Checkout',
      callbackConfirm: () => checkoutAction(hash),
    });
  }
  else {
    checkoutAction(hash);
  }
}

function checkoutAction(hash: string) {
  CheckoutCommit(hash).then(result => {
    parseResponse(result, () => {
      commitData.refresh();
      commitSignData.refresh();
      currentBranch.detach();
    });
  });
}

export default checkoutCommit;
