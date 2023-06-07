<script lang="ts">
  import { DiffConflictType, currentDiff } from "stores/diffs";

  let changesStatus: number[] = [];

  function hoverMiniHunk(minihunk: number, oursTheirs: number) {
    // select each time to capture updates
    let divs: NodeListOf<HTMLElement>;
    if (oursTheirs === DiffConflictType.Both) {
      divs = document.querySelectorAll(`.diff__row[data-minihunk="${minihunk}"]`);
    }
    else {
      divs = document.querySelectorAll(`.diff__row[data-minihunk="${minihunk}"][data-ourstheirs="${oursTheirs}"]`);
    }
    for (let i = 0; i < divs.length; i++) {
      if (parseInt(divs[i].dataset.minihunk) === minihunk) {
        divs[i].classList.add('diff__row--hover');
      } else {
        divs[i].classList.remove('diff__row--hover');
      }
    }
  }

  function unHoverMiniHunk(minihunk: number) {
    // select each time to capture updates
    let divs = document.querySelectorAll(`.diff__row[data-minihunk="${minihunk}"]`) as NodeListOf<HTMLElement>;
    for (let i = 0; i < divs.length; i++) {
      divs[i].classList.remove('diff__row--hover');
    }
  }

  function selectConflict(minihunk: number, oursTheirs: number) {
    currentDiff.setConflictResolution(minihunk, oursTheirs);
  }
</script>

<div class="diff diff--conflict">
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
          <div class="diff__row  diff__row--{line.Type}"
            data-minihunk="{line.MiniHunk}"
            data-ourstheirs="{line.OursTheirs}"
            class:diff__row--yes={
              $currentDiff.Conflicts[line.MiniHunk] && (
                $currentDiff.Conflicts[line.MiniHunk].Resolution === line.OursTheirs
                || $currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.Both
                || $currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.BothInverse
              )
            }
            class:diff__row--no={
              $currentDiff.Conflicts[line.MiniHunk]?.Resolution
              && $currentDiff.Conflicts[line.MiniHunk].Resolution !== line.OursTheirs
              && $currentDiff.Conflicts[line.MiniHunk].Resolution !== DiffConflictType.Both
              && $currentDiff.Conflicts[line.MiniHunk].Resolution !== DiffConflictType.BothInverse
            }
          >
            {#if line.Type === 'DiffChangeStartLine'}
              <div class="diff__change-start">
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Ours)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Ours)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflictType.Ours)}
                  class:btn--active={$currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.Ours}
                >Accept Current</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Theirs)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Theirs)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflictType.Theirs)}
                  class:btn--active={$currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.Theirs}
                >Accept Incoming</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Both)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Both)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflictType.Both)}
                  class:btn--active={$currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.Both}
                >Accept Both (Current First)</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Both)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflictType.Both)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflictType.BothInverse)}
                  class:btn--active={$currentDiff.Conflicts[line.MiniHunk].Resolution === DiffConflictType.BothInverse}
                >Accept Both (Incoming First)</button>
              </div>
            {/if}
            <div class="diff__line-toggle-minihunk"></div>
            <div class="diff__line diff__line--{line.Type} diff__line--noclick">
              <div class="diff__line-no">{
                line.Type === 'DiffDeleteLine'
                  ? line.OldLineNo
                : line.Type === 'DiffAddLine' || line.Type === 'DiffContextLine'
                  ? line.NewLineNo
                : ''
              }</div>
              <div class="diff__line-type"></div>
              <div class="diff__line-code" class:diff__line-code--nonewline={line.NoNewline}>{line.Line}</div>
            </div>
          </div>
        {/each}
      {/each}
    </div>
  {/if}
</div>
