<script lang="ts">
  import octicons from "@primer/octicons";
  import { currentDiff, isConflict } from "stores/diffs";

  let applyBtn: HTMLButtonElement;

  function applyChanges() {
    currentDiff.setResolved(true);
    applyBtn.disabled = true;
  }
</script>

<div class="diff-actions">
  {#if isConflict($currentDiff)}
    <button
      class="btn btn--apply"
      disabled={$currentDiff.NumConflicts !== Object.keys($currentDiff.ConflictSelections).length}
      on:click={applyChanges}
      bind:this={applyBtn}
    >
      Apply Selected Conflict Resolutions
      {@html octicons['check'].toSVG({width: 16})}
    </button>
  {/if}
</div>
