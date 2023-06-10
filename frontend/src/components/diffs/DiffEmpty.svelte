<script lang="ts">
  import { currentDiff } from "stores/diffs";
</script>

<div class="diff-message diff-message--empty">
  {#if $currentDiff.Flags.includes('replace-line-endings')}
    {#if $currentDiff.Committed}
      Line endings were updated in this file.
    {:else}
      Line endings will be updated in this file.
    {/if}
  {:else}
    Empty file.
  {/if}
</div>
