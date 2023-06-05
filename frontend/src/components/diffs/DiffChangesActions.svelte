<script lang="ts">
  import octicons from "@primer/octicons";
  import { parseResponse } from "scripts/parse-response";
  import { changes } from "stores/changes";
  import { currentDiff } from "stores/diffs";
  import { StageFile, StagePartialFile, UnstageFile, UnstagePartialFile } from "wailsjs/go/main/App";

  function stageSelected(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLButtonElement }) {
    e.currentTarget.disabled = true;
    StagePartialFile($currentDiff, $currentDiff.File, $currentDiff.Status).then(result => {
      parseResponse(result, () => {
        changes.refresh();
        currentDiff.refresh();
      }, () => {
        e.currentTarget.disabled = false;
      });
    });
  }

  function stageAll(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLButtonElement }) {
    e.currentTarget.disabled = true;
    StageFile($currentDiff.File).then(result => {
      parseResponse(result, () => {
        changes.refresh();
        currentDiff.clear();
      }, () => {
        e.currentTarget.disabled = false;
      });
    });
  }

  function unstageSelected(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLButtonElement }) {
    e.currentTarget.disabled = true;
    UnstagePartialFile($currentDiff, $currentDiff.File, $currentDiff.Status).then(result => {
      parseResponse(result, () => {
        changes.refresh();
        currentDiff.refresh();
      }, () => {
        e.currentTarget.disabled = false;
      });
    });
  }

  function unstageAll(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLButtonElement }) {
    e.currentTarget.disabled = true;
    UnstageFile($currentDiff.File).then(result => {
      parseResponse(result, () => {
        changes.refresh();
        currentDiff.clear();
      }, () => {
        e.currentTarget.disabled = false;
      });
    });
  }

  function setAllLines(selected: boolean) {
    for (let i = 0; i < $currentDiff.Hunks.length; i++) {
      for (let j = 0; j < $currentDiff.Hunks[i].Lines.length; j++) {
        $currentDiff.Hunks[i].Lines[j].Selected = selected;
      }
    }
    $currentDiff.SelectedLines = selected ? $currentDiff.SelectableLines : 0;
  }
</script>

<div class="diff-actions">
  {#if $currentDiff.Staged}
    <button
      class="btn"
      disabled={!$currentDiff.SelectedLines}
      on:click={$currentDiff.SelectedLines < $currentDiff.SelectableLines ? unstageSelected : unstageAll}
    >
      Unstage Selected Lines
      {@html octicons['arrow-down'].toSVG({width: 16})}
    </button>
  {:else}
    <button
      class="btn"
      disabled={!$currentDiff.SelectedLines}
      on:click={$currentDiff.SelectedLines < $currentDiff.SelectableLines ? stageSelected : stageAll}
    >
      Stage Selected Lines
      {@html octicons['arrow-up'].toSVG({width: 16})}
    </button>
  {/if}
  <button
    class="btn"
    on:click={() => setAllLines(true)}
    disabled={$currentDiff.SelectedLines === $currentDiff.SelectableLines}
  >
    Select All
    {@html octicons['check'].toSVG({width: 16})}
  </button>
  <button
    class="btn"
    on:click={() => setAllLines(false)}
    disabled={!$currentDiff.SelectedLines}
  >
    Deselect All
    {@html octicons['x'].toSVG({width: 16})}
  </button>
</div>
