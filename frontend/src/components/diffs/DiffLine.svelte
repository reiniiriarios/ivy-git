<script lang="ts">
  import { currentDiff, type DiffLine } from "stores/diffs";

  export let line: DiffLine;
  export let hunkId: number = -1;

  let selected: string;
  $: selected = !$currentDiff.Committed && !$currentDiff.Conflict
    ? (line.Selected ? 'diff__line--on' : 'diff__line--off')
    : '';
</script>

<div
  class="diff__line diff__line--{line.Type} {selected}"
  class:diff__line--noclick={$currentDiff.Committed || $currentDiff.Conflict}
  data-hunk="{hunkId}"
  data-minihunk="{line.MiniHunk}"
  on:click
  on:keypress
>
  {#if !$currentDiff.Committed && !$currentDiff.Conflict}
    <div class="diff__line-toggle"></div>
  {/if}
  <div class="diff__line-no">{
    line.Type === 'DiffDeleteLine'
      ? line.OldLineNo
    : line.Type === 'DiffAddLine' || line.Type === 'DiffContextLine'
      ? line.NewLineNo
    : ''
  }</div>
  <div class="diff__line-type"></div>
  <div class="diff__line-code" class:diff__line-code--nonewline={line.NoNewline}>
    <span class="diff__line-code-contents highlight">
      {#if $currentDiff.Highlight && $currentDiff.Highlight[line.CurLineNo]}
        {@html $currentDiff.Highlight[line.CurLineNo]}
      {:else if $currentDiff.Highlight}
        <span class="mute">{line.Line}</span>
      {:else}
        {line.Line}
      {/if}
    </span>
  </div>
</div>
