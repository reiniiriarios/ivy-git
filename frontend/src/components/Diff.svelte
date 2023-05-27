<script lang="ts">
  import { unstagedFileDiff } from "stores/diffs";

  function toggleLine(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    e.currentTarget.classList.toggle('diff__line--on');
    e.currentTarget.classList.toggle('diff__line--off');
  }

  let miniHunkElements: HTMLElement[] = [];
  let lineElements: HTMLElement[] = [];

  function hoverMiniHunk(e: (MouseEvent | FocusEvent) & { currentTarget: HTMLElement }) {
    miniHunkElements.map(line => {
      if (line.dataset.minihunk === e.currentTarget.dataset.minihunk) {
        line.classList.add('diff__line-toggle-minihunk--hover');
      } else {
        line.classList.remove('diff__line-toggle-minihunk--hover');
      }
    });
  }

  function unHoverMiniHunk(e: (MouseEvent | FocusEvent) & { currentTarget: HTMLElement }) {
    miniHunkElements.map(line => line.classList.remove('diff__line-toggle-minihunk--hover'));
  }

  function toggleMiniHunk(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    // Whether more lines are on or off.
    let onOff = lineElements.reduce(([on, off], ln) => {
      if (ln.dataset.minihunk === e.currentTarget.dataset.minihunk) {
        if (ln.classList.contains('diff__line--on')) {
          return [on + 1, off];
        } else {
          return [on, off + 1];
        }
      }
      return [on, off];
    }, [0, 0]);
    let moreLinesOn = onOff[0] > onOff[1];
    console.log(moreLinesOn)

    lineElements.forEach(line => {
      if (line.dataset.minihunk === e.currentTarget.dataset.minihunk) {
        if (moreLinesOn && line.classList.contains('diff__line--on')) {
          line.classList.remove('diff__line--on');
          line.classList.add('diff__line--off');
        }
        else if (!moreLinesOn && line.classList.contains('diff__line--off')) {
          line.classList.add('diff__line--on');
          line.classList.remove('diff__line--off');
        }
      }
    });
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
            <div class="diff__row">
              <div class="diff__line-toggle-minihunk"
                data-minihunk="{line.MiniHunk}"
                on:mouseover={hoverMiniHunk}
                on:mouseout={unHoverMiniHunk}
                on:focus={hoverMiniHunk}
                on:blur={unHoverMiniHunk}
                on:click={toggleMiniHunk}
                on:keypress={toggleMiniHunk}
                bind:this={miniHunkElements[line.RawLineNo]}
              ></div>
              <div class="diff__line diff__line--{line.Type} diff__line--on"
                data-minihunk="{line.MiniHunk}"
                on:click={toggleLine}
                on:keypress={toggleLine}
                bind:this={lineElements[line.RawLineNo]}
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
