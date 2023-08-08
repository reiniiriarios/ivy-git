import { derived, writable } from "svelte/store";
import { GetRepoState } from "wailsjs/go/ivy/App";
import { commitMessage } from "stores/ui";

export const RepoState = {
  Nil: "", // Default, no current repo, etc.
  None: "RepoStateNone",
  Interactive: "RepoStateRebaseInteractive",
  RebaseMerge: "RepoStateRebaseMerge",
  Merge: "RepoStateMerge",
  Rebase: "RepoStateRebase",
  Apply: "RepoStateApply",
  ApplyOrRebase: "RepoStateApplyOrRebase",
  Revert: "RepoStateRevert",
  RevertSequence: "RepoStateRevertSequence",
  CherryPick: "RepoStateCherryPick",
  CherryPickSequence: "RepoStateCherryPickSequence",
  Bisect: "RepoStateBisect",
}

function createRepoState() {
  const { subscribe, set } = writable(RepoState.Nil);
  
  return {
    subscribe,
    load: async () => {
      GetRepoState().then(result => {
        set(result);
        commitMessage.refresh();
      });
    },
    refresh: async () => {
      GetRepoState().then(result => {
        set(result);
        commitMessage.check();
      });
    },
    clear: async () => set(RepoState.Nil),
  };
}
export const repoState = createRepoState();
export const repoStateMessage = derived(repoState, $repoState => {
  switch ($repoState) {
    case RepoState.Interactive:
      return "Rebase in Progress"; // Interactive not always correct.
    case RepoState.RebaseMerge:
      return "Rebase or Merge in Progress";
    case RepoState.Merge:
      return "Merge in Progress";
    case RepoState.Rebase:
      return "Rebase in Progress";
    case RepoState.Apply:
      return "Apply in Progress";
    case RepoState.ApplyOrRebase:
      return "Apply or Rebase in Progress";
    case RepoState.Revert:
      return "Revert in Progress";
    case RepoState.RevertSequence:
      return "Revert Sequence in Progress";
    case RepoState.CherryPick:
      return "Cherry Pick in Progress";
    case RepoState.CherryPickSequence:
      return "Cherry Pick Sequence in Progress";
    case RepoState.Bisect:
      return "Bisect in Progress";
    default:
      return "";
  }
});
