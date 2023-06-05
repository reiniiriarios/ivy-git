<script lang="ts">
  import octicons from "@primer/octicons";
  import { currentDiff } from "stores/diffs";

  let applyBtn: HTMLButtonElement;

  function applyChanges() {
    currentDiff.setResolved(true);
    applyBtn.disabled = true;
  }
</script>

<div class="diff-actions diff-actions--conflict">
  <button
    class="btn btn--apply"
    disabled={$currentDiff.NumConflicts !== Object.keys($currentDiff.ConflictSelections).length}
    on:click={applyChanges}
    bind:this={applyBtn}
  >
    Apply Selected Conflict Resolutions
    {@html octicons['check'].toSVG({width: 16})}
  </button>
</div>
