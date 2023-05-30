<script lang="ts">
  import { currentDiff, type DiffLine } from "stores/diffs";

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
    let onOff = $currentDiff.Hunks[hunk].Lines.reduce(([on, off], ln: DiffLine) => {
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

    for (let i = 0; i < $currentDiff.Hunks[hunk].Lines.length; i++) {
      if ($currentDiff.Hunks[hunk].Lines[i].MiniHunk === miniHunk) {
        $currentDiff.Hunks[hunk].Lines[i].Selected = moreLinesOn ? false : true;
      }
    }
  }
</script>

<div class="diff">
  {#if $currentDiff.Hunks?.length}
    <div class="diff__grid">
      {#each $currentDiff.Hunks as hunk, h}
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
            <div class="diff__row">
              <div class="diff__line-toggle-minihunk diff__line-toggle-minihunk--click"
                data-hunk="{h}"
                data-minihunk="{line.MiniHunk}"
                on:mouseover={hoverMiniHunk}
                on:mouseout={unHoverMiniHunk}
                on:focus={hoverMiniHunk}
                on:blur={unHoverMiniHunk}
                on:click={() => toggleMiniHunk(h, line.MiniHunk)}
                on:keypress={() => toggleMiniHunk(h, line.MiniHunk)}
                bind:this={miniHunkElements[line.RawLineNo]}
              ></div>
              <div class="diff__line diff__line--{line.Type} diff__line--{line.Selected ? 'on' : 'off'}"
                data-hunk="{h}"
                data-minihunk="{line.MiniHunk}"
                on:click={() => line.Selected = !line.Selected}
                on:keypress={() => line.Selected = !line.Selected}
              >
                <div class="diff__line-toggle"></div>
                <div class="diff__line-no">{line.Type === 'DiffDeleteLine' ? line.OldLineNo : line.NewLineNo}</div>
                <div class="diff__line-type"></div>
                <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
              </div>
            </div>
          {:else}
            <div class="diff__row">
              <div class="diff__line-toggle-minihunk"></div>
              <div class="diff__line diff__line--{line.Type}">
                <div class="diff__line-toggle"></div>
                <div class="diff__line-no">{line.NewLineNo}</div>
                <div class="diff__line-type"></div>
                <div class="diff__line-code" class:diff__line--nonewline={line.NoNewline}>{line.Line}</div>
              </div>
            </div>
          {/if}
        {/each}
      {/each}
    </div>
  {/if}
</div>
