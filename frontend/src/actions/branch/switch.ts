import { currentBranch } from "stores/branches";
import { messageDialog } from "stores/message-dialog";
import { RepoState, repoState } from "stores/repo-state";
import { branchSelect } from "stores/ui";
import { get } from "svelte/store";

export default function switchBranch(b: string, r: string = '') {
  if (![RepoState.Nil, RepoState.None].includes(get(repoState))) {
    messageDialog.error({
      message: "The repo is currently in a state that you cannot (or should not) switch branches."
    });
  } else {
    currentBranch.switch(b, r);
    branchSelect.set(false);
  }
}
