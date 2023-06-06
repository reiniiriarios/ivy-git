<script lang="ts">
  import octicons from "@primer/octicons";
  import { currentDiff, currentDiffResolved } from "stores/diffs";

  let applyBtn: HTMLButtonElement;

  function applyChanges() {
    applyBtn.disabled = true;
    currentDiff.resolveConflicts();
  }
</script>

<div class="diff-actions diff-actions--conflict">
  <button
    class="btn btn--apply"
    disabled={!$currentDiffResolved}
    on:click={applyChanges}
    bind:this={applyBtn}
  >
    Apply Selected Conflict Resolutions
    {@html octicons['check'].toSVG({width: 16})}
  </button>
</div>
