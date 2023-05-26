<script lang="ts">
  import { unstagedFileDiff } from "stores/diffs";
</script>

<div class="diff">
  {#if $unstagedFileDiff.Hunks?.length}
    <div class="diff__grid">
      {#each $unstagedFileDiff.Hunks as hunk}
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
          <div class="diff__line-no diff-line--{line.Type}">{line.Type === 'DiffDeleteLine' ? line.OldLineNo : line.NewLineNo}</div>
          <div class="diff__line-type diff-line--{line.Type}"></div>
          <div class="diff__line diff-line--{line.Type}" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
        {/each}
      {/each}
    </div>
  {/if}
</div>
