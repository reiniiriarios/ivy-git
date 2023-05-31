<script lang="ts">
  import { currentDiff } from "stores/diffs";
  import ChangesDiff from "components/ChangesDiff.svelte";
  import ChangesDiffActions from "components/ChangesDiffActions.svelte";
  import DiffCommitted from "./DiffCommitted.svelte";
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
    {:else}
      {#if $currentDiff.File}
        <ChangesDiffActions />
      {/if}
      <ChangesDiff />
    {/if}
  </div>
{/if}
