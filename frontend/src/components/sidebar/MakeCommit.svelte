<script lang="ts">
  import { isDarwin } from "scripts/env";
  import { parseResponse } from "scripts/parse-response";
  import { currentBranch, detachedHead } from "stores/branches";
  import { changes, changesNumConflicts, changesNumStaged, changesNumUnstaged, mergeConflicts, mergeConflictsResolved } from "stores/changes";
  import { repoState } from "stores/repo-state";
  import { repoSelect, branchSelect, commitMessageSubject, commitMessageBody, commitMessage } from "stores/ui";
  import { MakeCommit } from "wailsjs/go/main/App";

  let running: boolean = false;

  function make() {
    if ($commitMessageSubject && !running && ($changesNumStaged || $changesNumUnstaged)) {
      running = true;
      MakeCommit($commitMessageSubject, $commitMessageBody).then(result => {
        parseResponse(result, () => {
          commitMessage.clear();
        });
        running = false;
      })
    }
  }

  function makeCommitKeypress(e: KeyboardEvent & { currentTarget: HTMLElement }) {
    let cmd = (isDarwin() && e.metaKey) || (!isDarwin() && e.ctrlKey);
    if (cmd && (e.key === '\n' || e.key === 'Enter')) {
      make();
      e.currentTarget.blur();
    }
  }
</script>

<div
  class="make-commit"
  class:detached={$detachedHead}
  style:display={$repoSelect || $branchSelect ? 'none' : 'block'}
>
  <div class="make-commit__subject">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Summary</span>
      <input class="text-input repo-state--{$repoState}" type="text" bind:value={$commitMessageSubject} on:keypress={makeCommitKeypress}>
    </label>
  </div>
  <div class="make-commit__body">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Description</span>
      <textarea class="text-input repo-state--{$repoState}" bind:value={$commitMessageBody} on:keypress={makeCommitKeypress}></textarea>
    </label>
  </div>
  <div class="make-commit__button">
    <button
      class="btn repo-state--{$repoState}"
      id="make-commit-button"
      disabled={
        !$commitMessageSubject
        || running
        || (!$changesNumStaged && !$changesNumUnstaged && !$changesNumConflicts)
        || ($mergeConflicts && !$mergeConflictsResolved)
      }
      on:click={make}
    >
      Commit to <strong>{$currentBranch?.Name}</strong>
    </button>
  </div>
</div>
