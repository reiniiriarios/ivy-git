<script lang="ts">
  import octicons from "@primer/octicons";
  import { currentDiff } from "stores/diffs";

  function stageSelected() {

  }

  function stageAll() {

  }

  function unstageSelected() {

  }

  function unstageAll() {

  }
</script>

<div class="diff-actions">
  {#if $currentDiff.Staged}
    {#if !$currentDiff.SelectedLines}
      <button class="btn" disabled>
        No Lines Selected
      </button>
    {:else if $currentDiff.SelectedLines < $currentDiff.SelectableLines}
      <button class="btn" on:click={unstageSelected}>
        Unstage Selected Lines
        {@html octicons['arrow-down'].toSVG({width: 16})}
      </button>
    {:else}
      <button class="btn" on:click={unstageAll}>
        Unstage File
        {@html octicons['arrow-down'].toSVG({width: 16})}
      </button>
    {/if}
  {:else}
    {#if !$currentDiff.SelectedLines}
      <button class="btn" disabled>
        No Lines Selected
      </button>
    {:else if $currentDiff.SelectedLines < $currentDiff.SelectableLines}
      <button class="btn" on:click={stageSelected}>
        Stage Selected Lines
        {@html octicons['arrow-up'].toSVG({width: 16})}
      </button>
    {:else}
      <button class="btn" on:click={stageAll}>
        Stage File
        {@html octicons['arrow-up'].toSVG({width: 16})}
      </button>
    {/if}
  {/if}
</div>
