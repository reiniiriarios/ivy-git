<script lang="ts">
  import { unstagedFileDiff } from "stores/diffs";

  function toggleLine(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    e.currentTarget.classList.toggle('diff__line--on');
    e.currentTarget.classList.toggle('diff__line--off');
  }
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
          {#if line.Type !== 'DiffContextLine'}
            <div class="diff__line diff__line--{line.Type} diff__line--on" on:click={toggleLine} on:keypress={toggleLine}>
              <div class="diff__line-toggle"></div>
              <div class="diff__line-no">{line.Type === 'DiffDeleteLine' ? line.OldLineNo : line.NewLineNo}</div>
              <div class="diff__line-type"></div>
              <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
            </div>
          {:else}
            <div class="diff__line diff__line--{line.Type}">
              <div class="diff__line-toggle"></div>
              <div class="diff__line-no">{line.NewLineNo}</div>
              <div class="diff__line-type"></div>
              <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
            </div>
          {/if}
        {/each}
      {/each}
    </div>
  {/if}
</div>
