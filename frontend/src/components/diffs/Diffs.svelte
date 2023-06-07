<script lang="ts">
  import { currentDiff } from "stores/diffs";
  import DiffChanges from "components/diffs/DiffChanges.svelte";
  import DiffChangesActions from "components/diffs/DiffChangesActions.svelte";
  import DiffCommitted from "components/diffs/DiffCommitted.svelte";
  import DiffConflict from "components/diffs/DiffConflict.svelte";
  import DiffConflictActions from "components/diffs/DiffConflictActions.svelte";
</script>

{#if $currentDiff}
  <div class="diffs">
    {#if $currentDiff.Committed}
      <div class="diffs__hash">
        {$currentDiff.Hash}
      </div>
      <div class="diffs__filename">
        {$currentDiff.File}
      </div>
      <DiffCommitted />
    {:else if $currentDiff.Conflict}
      <DiffConflictActions />
      <DiffConflict />
    {:else}
      {#if $currentDiff.File && !$currentDiff.Binary && !$currentDiff.Empty}
        <DiffChangesActions />
      {/if}
      <DiffChanges />
    {/if}
  </div>
{/if}
