<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { changes, mergeConflicts, mergeConflictsResolved } from "stores/changes";
  import { repoState } from "stores/repo-state";
  import { commitMessage } from "stores/ui";
  import { RebaseAbort, RebaseContinue, RebaseSkip } from "wailsjs/go/ivy/App";

  let running: boolean = false;

  function rebaseContinue() {
    if (!running) {
      running = true;
      RebaseContinue().then(result => {
        parseResponse(result, () => {
          changes.refresh();
          repoState.refresh();
          commitMessage.clear();
        });
        running = false;
      });
    }
  }

  function rebaseAbort() {
    if (!running) {
      running = true;
      RebaseAbort().then(result => {
        parseResponse(result, () => {
          changes.refresh();
          repoState.refresh();
          commitMessage.clear();
        });
        running = false;
      });
    }
  }

  function rebaseSkip() {
    if (!running) {
      running = true;
      RebaseSkip().then(result => {
        parseResponse(result, () => {
          changes.refresh();
          repoState.refresh();
          commitMessage.clear();
        });
        running = false;
      });
    }
  }
</script>

<div class="rebase-actions repo-state--{$repoState}">
  <div class="rebase-actions__name">Rebase In Progress</div>
  <div class="rebase-actions__actions">
    <div class="rebase-actions__action">
      <button
        class="btn btn--green"
        disabled={running || ($mergeConflicts && !$mergeConflictsResolved)}
        on:click={rebaseContinue}
      >Continue</button>
    </div>
    <div class="rebase-actions__action">
      <button
        class="btn"
        disabled={running}
        on:click={rebaseAbort}
      >Abort</button>
    </div>
    <div class="rebase-actions__action">
      <button
        class="btn"
        disabled={running}
        on:click={rebaseSkip}
      >Skip</button>
    </div>
  </div>
</div>
