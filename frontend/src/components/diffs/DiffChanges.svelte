<script lang="ts">
  import { changes, changesNumConflicts, changesNumStaged, changesNumUnstaged } from "stores/changes";
  import { currentDiff, type DiffLine as DiffLineData } from "stores/diffs";
  import DiffBinary from "components/diffs/DiffBinary.svelte";
  import DiffEmpty from "components/diffs/DiffEmpty.svelte";
  import NoChanges from "./NoChanges.svelte";
  import DiffLine from "./DiffLine.svelte";
  import ChangesSummary from "./ChangesSummary.svelte";

  let miniHunkElements: HTMLElement[] = [];

  function hoverMiniHunk(e: (MouseEvent | FocusEvent) & { currentTarget: HTMLElement }) {
    // select each time to capture updates
    let divs = document.querySelectorAll(`.diff__line-toggle-minihunk[data-minihunk="${e.currentTarget.dataset.minihunk}"]`) as NodeListOf<HTMLElement>;
    for (let i = 0; i < divs.length; i++) {
      if (divs[i].dataset.minihunk === e.currentTarget.dataset.minihunk) {
        divs[i].classList.add('diff__line-toggle-minihunk--hover');
      } else {
        divs[i].classList.remove('diff__line-toggle-minihunk--hover');
      }
    }
  }

  function unHoverMiniHunk(e: (MouseEvent | FocusEvent) & { currentTarget: HTMLElement }) {
    // select each time to capture updates
    let divs = document.querySelectorAll(`.diff__line-toggle-minihunk[data-minihunk="${e.currentTarget.dataset.minihunk}"]`) as NodeListOf<HTMLElement>;
    for (let i = 0; i < divs.length; i++) {
      divs[i].classList.remove('diff__line-toggle-minihunk--hover');
    }
  }

  function toggleMiniHunk(hunk: number, miniHunk: number) {
    // Whether more lines are on or off.
    let onOff = $currentDiff.Hunks[hunk].Lines.reduce(([on, off], ln: DiffLineData) => {
      if (ln.MiniHunk === miniHunk) {
        if (ln.Selected) {
          return [on + 1, off];
        } else {
          return [on, off + 1];
        }
      }
      return [on, off];
    }, [0, 0]);
    let moreLinesOn = onOff[0] > onOff[1];

    // Toggle minihunk lines all on or all off.
    // Adjust number of selected lines for diff.
    let adj = 0;
    for (let i = 0; i < $currentDiff.Hunks[hunk].Lines.length; i++) {
      if ($currentDiff.Hunks[hunk].Lines[i].MiniHunk === miniHunk) {
        if (moreLinesOn && $currentDiff.Hunks[hunk].Lines[i].Selected) {
          adj--;
        } else if (!moreLinesOn && !$currentDiff.Hunks[hunk].Lines[i].Selected) {
          adj++;
        }
        $currentDiff.Hunks[hunk].Lines[i].Selected = !moreLinesOn;
      }
    }
    $currentDiff.SelectedLines += adj;
    changes.setPartialFromCurrent();
  }

  function toggleLine(hunk: number, i: number) {
    $currentDiff.Hunks[hunk].Lines[i].Selected = !$currentDiff.Hunks[hunk].Lines[i].Selected;
    $currentDiff.SelectedLines += $currentDiff.Hunks[hunk].Lines[i].Selected ? 1 : -1;
    changes.setPartialFromCurrent();
  }
</script>

<div class="diff diff--changes">
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
        {#each hunk.Lines as line, line_id}
          {#if line.Type !== 'DiffContextLine'}
            <div class="diff__row">
              <div class="diff__line-toggle-minihunk diff__line-toggle-minihunk--click"
                data-hunk="{hunk_id}"
                data-minihunk="{line.MiniHunk}"
                on:mouseover={hoverMiniHunk}
                on:mouseout={unHoverMiniHunk}
                on:focus={hoverMiniHunk}
                on:blur={unHoverMiniHunk}
                on:click={() => toggleMiniHunk(hunk_id, line.MiniHunk)}
                on:keypress={() => toggleMiniHunk(hunk_id, line.MiniHunk)}
                bind:this={miniHunkElements[line.RawLineNo]}
              ></div>
              <DiffLine bind:line={line} hunkId={hunk_id} 
                on:click={() => toggleLine(hunk_id, line_id)}
                on:keypress={() => toggleLine(hunk_id, line_id)}
              />
            </div>
          {:else}
            <div class="diff__row">
              <div class="diff__line-toggle-minihunk"></div>
              <DiffLine bind:line={line} />
            </div>
          {/if}
        {/each}
      {/each}
    </div>
  {:else if $currentDiff.Binary}
    <DiffBinary />
  {:else if $currentDiff.Empty}
    <DiffEmpty />
  {:else if $changesNumStaged || $changesNumUnstaged || $changesNumConflicts}
    <ChangesSummary />
  {:else}
    <NoChanges />
  {/if}
</div>
