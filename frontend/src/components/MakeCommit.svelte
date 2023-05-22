<script lang="ts">
  import { currentBranch } from "stores/branches";
  import { repoSelect, branchSelect } from "stores/ui";

  let subject: string;
  let body: string;

  function make() {
    if (subject) {
      //...
    }
  }
</script>

<div class="make-commit" style:display={$repoSelect || $branchSelect ? 'none' : 'block'}>
  <div class="make-commit__subject">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Summary</span>
      <input type="text" bind:value={subject}>
    </label>
  </div>
  <div class="make-commit__body">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Description</span>
      <textarea bind:value={body}></textarea>
    </label>
  </div>
  <div class="make-commit__button">
    <button class="btn" id="make-commit-button" disabled={!subject} on:click={make}>
      Commit to <strong>{$currentBranch.Name}</strong>
    </button>
  </div>
</div>
