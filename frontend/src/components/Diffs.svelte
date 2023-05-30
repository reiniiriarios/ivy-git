<script lang="ts">
  import { currentDiff } from "stores/diffs";
  import ChangesDiff from "components/ChangesDiff.svelte";
  import ChangesDiffActions from "components/ChangesDiffActions.svelte";
</script>

{#if $currentDiff}
  <div class="diffs">
    {#if $currentDiff.Committed}
      todo
    {:else}
      <ChangesDiffActions />
      <ChangesDiff />
    {/if}
  </div>
{/if}
