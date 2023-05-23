<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { currentBranch } from "stores/branches";
  import { repoSelect, branchSelect } from "stores/ui";
  import { MakeCommit } from "wailsjs/go/main/App";

  let subject: string;
  let body: string;
  let running: boolean = false;

  function make() {
    if (subject) {
      running = true;
      MakeCommit(subject, body).then(result => {
        parseResponse(result, () => {
          subject = '';
          body = '';
        });
        running = false;
      })
    }
  }

  function makeCommitKeypress(e: KeyboardEvent & { currentTarget: HTMLElement }) {
    if ((e.metaKey || e.ctrlKey) && e.key === '\n') {
      make();
      e.currentTarget.blur();
    }
  }
</script>

<div class="make-commit" style:display={$repoSelect || $branchSelect ? 'none' : 'block'}>
  <div class="make-commit__subject">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Summary</span>
      <input type="text" bind:value={subject} on:keypress={makeCommitKeypress}>
    </label>
  </div>
  <div class="make-commit__body">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Description</span>
      <textarea bind:value={body} on:keypress={makeCommitKeypress}></textarea>
    </label>
  </div>
  <div class="make-commit__button">
    <button class="btn" id="make-commit-button" disabled={!subject || running} on:click={make}>
      Commit to <strong>{$currentBranch.Name}</strong>
    </button>
  </div>
</div>
