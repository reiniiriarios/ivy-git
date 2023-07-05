import { currentBranch, detachedHead } from "stores/branches";
import { messageDialog } from "stores/message-dialog";
import { RepoState, repoState } from "stores/repo-state";
import { branchSelect } from "stores/ui";
import { get } from "svelte/store";

export default function switchBranch(branch: string, remote: string = '') {
  if (![RepoState.Nil, RepoState.None].includes(get(repoState))) {
    messageDialog.error({
      message: "The repo is currently in a state that you cannot (or should not) switch branches."
      // add confirm to checkout anyway??
    });
  }
  else if (get(detachedHead)) {
    messageDialog.confirm({
      heading: 'Checkout Branch',
      message: 'You are currently in a <strong>detached HEAD</strong> state. Checking out a branch could result in lost work. Continue?',
      confirm: 'Checkout',
      callbackConfirm: () => currentBranch.switch(branch, remote),
    });
  }
  else {
    currentBranch.switch(branch, remote);
    branchSelect.set(false);
  }
}
