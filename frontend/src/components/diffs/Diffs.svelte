<script lang="ts">
  import { currentDiff, isConflict, isDiff } from "stores/diffs";
  import DiffChanges from "components/diffs/DiffChanges.svelte";
  import DiffChangesActions from "components/diffs/DiffChangesActions.svelte";
  import DiffCommitted from "components/diffs/DiffCommitted.svelte";
  import DiffConflict from "components/diffs/DiffConflict.svelte";
</script>

{#if $currentDiff}
  <div class="diffs">
    {#if isDiff($currentDiff) && $currentDiff.Committed}
      <div class="diffs__hash">
        {$currentDiff.Hash}
      </div>
      <div class="diffs__filename">
        {$currentDiff.File}
      </div>
      <DiffCommitted />
    {:else if isConflict($currentDiff)}
      <DiffConflict />
    {:else}
      {#if $currentDiff.File}
        <DiffChangesActions />
      {/if}
      <DiffChanges />
    {/if}
  </div>
{/if}
