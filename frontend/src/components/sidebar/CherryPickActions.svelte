<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { changes, mergeConflicts, mergeConflictsResolved } from "stores/changes";
  import { repoState } from "stores/repo-state";
  import { commitMessage } from "stores/ui";
  import { CherryPickAbort, CherryPickContinue, CherryPickSkip } from "wailsjs/go/ivy/App";

  let running: boolean = false;

  function cherryContinue() {
    if (!running) {
      running = true;
      CherryPickContinue().then(result => {
        parseResponse(result, () => {
          changes.refresh();
          repoState.refresh();
          commitMessage.clear();
        });
        running = false;
      });
    }
  }

  function cherryAbort() {
    if (!running) {
      running = true;
      CherryPickAbort().then(result => {
        parseResponse(result, () => {
          changes.refresh();
          repoState.refresh();
          commitMessage.clear();
        });
        running = false;
      });
    }
  }

  function cherrySkip() {
    if (!running) {
      running = true;
      CherryPickSkip().then(result => {
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

<div class="cherry-pick-actions repo-state--{$repoState}">
  <div class="cherry-pick-actions__name">Cherry Pick In Progress</div>
  <div class="cherry-pick-actions__actions">
    <div class="cherry-pick-actions__action">
      <button
        class="btn btn--green"
        disabled={running || ($mergeConflicts && !$mergeConflictsResolved)}
        on:click={cherryContinue}
      >Continue</button>
    </div>
    <div class="cherry-pick-actions__action">
      <button
        class="btn"
        disabled={running}
        on:click={cherryAbort}
      >Abort</button>
    </div>
    <div class="cherry-pick-actions__action">
      <button
        class="btn"
        disabled={running}
        on:click={cherrySkip}
      >Skip</button>
    </div>
  </div>
</div>
