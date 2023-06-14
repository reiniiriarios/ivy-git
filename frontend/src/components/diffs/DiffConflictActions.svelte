<script lang="ts">
  import octicons from "@primer/octicons";
  import { changes } from "stores/changes";
  import { currentDiff, currentDiffResolved } from "stores/diffs";

  let applyBtn: HTMLButtonElement;
  let deleteBtn: HTMLButtonElement;
  let keepBtn: HTMLButtonElement;

  function applyChanges() {
    applyBtn.disabled = true;
    currentDiff.resolveConflicts();
  }

  function deleteFile() {
    deleteBtn.disabled = true;
    currentDiff.delete();
  }

  function keepFile() {
    keepBtn.disabled = true;
    currentDiff.keep();
  }
</script>

<div class="diff-actions diff-actions--conflict">
  {#if $changes.c[$currentDiff.File] && ($changes.c[$currentDiff.File].Them === 'D' || $changes.c[$currentDiff.File].Us === 'D')}
    <button
      class="btn btn--delete"
      on:click={deleteFile}
      bind:this={deleteBtn}
    >
      Delete File
      {@html octicons['file-removed'].toSVG({width: 14})}
    </button>
    {#if $currentDiff.Conflicts?.length}
      <button
        class="btn btn--apply"
        disabled={!$currentDiffResolved}
        on:click={applyChanges}
        bind:this={applyBtn}
      >
        Apply Selected Conflict Resolutions
        {@html octicons['check'].toSVG({width: 16})}
      </button>
    {:else}
      <button
        class="btn btn--apply"
        disabled={!$currentDiffResolved}
        on:click={keepFile}
        bind:this={keepBtn}
      >
        Keep File
        {@html octicons['file-added'].toSVG({width: 14})}
      </button>
    {/if}
  {:else if $currentDiff.Conflicts?.length}
    <button
      class="btn btn--apply"
      disabled={!$currentDiffResolved}
      on:click={applyChanges}
      bind:this={applyBtn}
    >
      Apply Selected Conflict Resolutions
      {@html octicons['check'].toSVG({width: 16})}
    </button>
  {/if}
</div>
