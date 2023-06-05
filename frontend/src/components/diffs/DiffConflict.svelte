<script lang="ts">
  import { DiffConflict, currentDiff, isConflict } from "stores/diffs";

  function hoverMiniHunk(conflict: number, oursTheirs: number) {
    // select each time to capture updates
    let divs: NodeListOf<HTMLElement>;
    if (oursTheirs === DiffConflict.Both) {
      divs = document.querySelectorAll(`.diff__row[data-conflict="${conflict}"]`);
    }
    else {
      divs = document.querySelectorAll(`.diff__row[data-conflict="${conflict}"][data-ourstheirs="${oursTheirs}"]`);
    }
    for (let i = 0; i < divs.length; i++) {
      if (parseInt(divs[i].dataset.conflict) === conflict) {
        divs[i].classList.add('diff__row--hover');
      } else {
        divs[i].classList.remove('diff__row--hover');
      }
    }
  }

  function unHoverMiniHunk(conflict: number) {
    // select each time to capture updates
    let divs = document.querySelectorAll(`.diff__row[data-conflict="${conflict}"]`) as NodeListOf<HTMLElement>;
    for (let i = 0; i < divs.length; i++) {
      divs[i].classList.remove('diff__row--hover');
    }
  }
</script>

<div class="diff diff--conflict">
  {#if isConflict($currentDiff) && $currentDiff.Lines?.length}
    <div class="diff__grid">
      {#each $currentDiff.Lines as line}
        <div class="diff__row  diff__row--{line.Type}"
          data-conflict="{line.Conflict}"
          data-ourstheirs="{line.OursTheirs}"
          class:diff__row--yes={
            $currentDiff.ConflictSelections[line.Conflict] === line.OursTheirs
            || $currentDiff.ConflictSelections[line.Conflict] >= 2
          }
          class:diff__row--no={
            $currentDiff.ConflictSelections[line.Conflict]
            && $currentDiff.ConflictSelections[line.Conflict] !== line.OursTheirs
            && $currentDiff.ConflictSelections[line.Conflict] < 2
          }
        >
          {#if line.Type === 'LineConflictStart'}
            <div class="diff__change-start">
              <button class="btn btn-text"
                on:mouseover={() => hoverMiniHunk(line.Conflict, DiffConflict.Ours)}
                on:mouseout={() => unHoverMiniHunk(line.Conflict)}
                on:focus={() => hoverMiniHunk(line.Conflict, DiffConflict.Ours)}
                on:blur={() => unHoverMiniHunk(line.Conflict)}
                on:click={() => currentDiff.setConflictResolution(line.Conflict, DiffConflict.Ours)}
                class:btn--active={$currentDiff.ConflictSelections[line.Conflict] === DiffConflict.Ours}
              >Accept Current</button>
              <button class="btn btn-text"
                on:mouseover={() => hoverMiniHunk(line.Conflict, DiffConflict.Theirs)}
                on:mouseout={() => unHoverMiniHunk(line.Conflict)}
                on:focus={() => hoverMiniHunk(line.Conflict, DiffConflict.Theirs)}
                on:blur={() => unHoverMiniHunk(line.Conflict)}
                on:click={() => currentDiff.setConflictResolution(line.Conflict, DiffConflict.Theirs)}
                class:btn--active={$currentDiff.ConflictSelections[line.Conflict] === DiffConflict.Theirs}
              >Accept Incoming</button>
              <button class="btn btn-text"
                on:mouseover={() => hoverMiniHunk(line.Conflict, DiffConflict.Both)}
                on:mouseout={() => unHoverMiniHunk(line.Conflict)}
                on:focus={() => hoverMiniHunk(line.Conflict, DiffConflict.Both)}
                on:blur={() => unHoverMiniHunk(line.Conflict)}
                on:click={() => currentDiff.setConflictResolution(line.Conflict, DiffConflict.Both)}
                class:btn--active={$currentDiff.ConflictSelections[line.Conflict] === DiffConflict.Both}
              >Accept Both (Current First)</button>
              <button class="btn btn-text"
                on:mouseover={() => hoverMiniHunk(line.Conflict, DiffConflict.Both)}
                on:mouseout={() => unHoverMiniHunk(line.Conflict)}
                on:focus={() => hoverMiniHunk(line.Conflict, DiffConflict.Both)}
                on:blur={() => unHoverMiniHunk(line.Conflict)}
                on:click={() => currentDiff.setConflictResolution(line.Conflict, DiffConflict.BothInverse)}
                class:btn--active={$currentDiff.ConflictSelections[line.Conflict] === DiffConflict.BothInverse}
              >Accept Both (Incoming First)</button>
            </div>
          {/if}
          <div class="diff__line-toggle-minihunk"></div>
          <div class="diff__line diff__line--{line.Type} diff__line--noclick">
            <div class="diff__line-no">{line.LineNo}</div>
            <div class="diff__line-code">{line.Line}</div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
