<script lang="ts">
  import { conflicts } from "stores/conflicts";
  import { DiffConflict, currentDiff } from "stores/diffs";

  let changesStatus: number[] = [];

  function hoverMiniHunk(minihunk: number, oursTheirs: number) {
    // select each time to capture updates
    let divs: NodeListOf<HTMLElement>;
    if (oursTheirs === DiffConflict.Both) {
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
    changesStatus[minihunk] = oursTheirs;
    $conflicts[$currentDiff.File].ConflictSelections[minihunk] = oursTheirs;
    changesStatus = $conflicts[$currentDiff.File].ConflictSelections;
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
            class:diff__row--yes={changesStatus[line.MiniHunk] === line.OursTheirs || changesStatus[line.MiniHunk] >= 2}
            class:diff__row--no={changesStatus[line.MiniHunk] && changesStatus[line.MiniHunk] !== line.OursTheirs && changesStatus[line.MiniHunk] < 2}
          >
            {#if line.Type === 'DiffChangeStartLine'}
              <div class="diff__change-start">
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Ours)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Ours)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflict.Ours)}
                  class:btn--active={changesStatus[line.MiniHunk] === DiffConflict.Ours}
                >Accept Current</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Theirs)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Theirs)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflict.Theirs)}
                  class:btn--active={changesStatus[line.MiniHunk] === DiffConflict.Theirs}
                >Accept Incoming</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Both)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Both)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflict.Both)}
                  class:btn--active={changesStatus[line.MiniHunk] === DiffConflict.Both}
                >Accept Both (Current First)</button>
                <button class="btn btn-text"
                  on:mouseover={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Both)}
                  on:mouseout={() => unHoverMiniHunk(line.MiniHunk)}
                  on:focus={() => hoverMiniHunk(line.MiniHunk, DiffConflict.Both)}
                  on:blur={() => unHoverMiniHunk(line.MiniHunk)}
                  on:click={() => selectConflict(line.MiniHunk, DiffConflict.BothInverse)}
                  class:btn--active={changesStatus[line.MiniHunk] === DiffConflict.BothInverse}
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
              <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
            </div>
          </div>
        {/each}
      {/each}
    </div>
  {/if}
</div>
