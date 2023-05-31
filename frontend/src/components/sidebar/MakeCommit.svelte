<script lang="ts">
  import { isDarwin } from "scripts/env";
  import { parseResponse } from "scripts/parse-response";
  import { currentBranch, detachedHead } from "stores/branches";
  import { changes } from "stores/changes";
  import { repoSelect, branchSelect, inProgressCommitMessage } from "stores/ui";
  import { MakeCommit } from "wailsjs/go/main/App";

  let subject: string;
  let body: string;

  inProgressCommitMessage.subscribe(msg => {
    subject = msg.Subject;
    body = msg.Body;
  });

  let running: boolean = false;

  function make() {
    if (subject && !running && (changes.numStaged() || changes.numUnstaged())) {
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
    <button class="btn" id="make-commit-button" disabled={!subject || running || (!changes.numStaged() && !changes.numUnstaged())} on:click={make}>
      Commit to <strong>{$currentBranch.Name}</strong>
    </button>
  </div>
</div>
