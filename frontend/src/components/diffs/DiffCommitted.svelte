<script lang="ts">
  import { currentDiff } from "stores/diffs";
  import DiffBinary from "components/diffs/DiffBinary.svelte";
  import DiffEmpty from "components/diffs/DiffEmpty.svelte";
  import DiffLine from "./DiffLine.svelte";
</script>

<div class="diff diff--committed">
  {#if $currentDiff.Hunks?.length}
    <div class="diff__grid">
      {#each $currentDiff.Hunks as hunk}
        <div class="diff__hunk-header">
          <span class="diff__hunk-details">
            @@
            -{hunk.StartOld},{hunk.EndOld}
            +{hunk.StartNew},{hunk.EndNew}
            @@
          </span>
          <span class="diff__hunk-heading">{hunk.Heading}</span>
        </div>
        {#each hunk.Lines as line}
          <div class="diff__row">
            <DiffLine line={line} />
          </div>
        {/each}
      {/each}
    </div>
  {:else if $currentDiff.Binary}
    <DiffBinary />
  {:else if $currentDiff.Empty}
    <DiffEmpty />
  {/if}
</div>
