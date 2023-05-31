<script lang="ts">
  import { currentDiff } from "stores/diffs";
</script>

<div class="diff">
  {#if $currentDiff.Hunks?.length}
    <div class="diff__grid">
      {#each $currentDiff.Hunks as hunk, hunk_id}
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
            <div class="diff__line-toggle-minihunk"></div>
            <div class="diff__line diff__line--{line.Type} diff__line--noclick">
              <div class="diff__line-toggle"></div>
              <div class="diff__line-no">{line.Type === 'DiffDeleteLine' ? line.OldLineNo : line.NewLineNo}</div>
              <div class="diff__line-type"></div>
              <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
            </div>
          </div>
        {/each}
      {/each}
    </div>
  {/if}
</div>
