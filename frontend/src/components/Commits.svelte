<script lang="ts">
  import CommitRow from 'components/CommitRow.svelte';
  import CommitDetails from 'components/CommitDetails.svelte';
  import { setCommitsContainer } from 'scripts/commit-details-resize';
  import { commitData, commits, tree, commitSignData } from 'stores/commit-data';
  import { afterUpdate, beforeUpdate, onMount } from 'svelte';
  import { setFade } from 'scripts/graph';
  import LoadMoreCommits from 'components/LoadMoreCommits.svelte';
  import { currentCommit } from 'stores/commit-details';
  import { get } from 'svelte/store';

  let scrollDiv: HTMLElement;
  let scrollPosition: number;
  let commitsTable: HTMLElement;
  let resizeIndex: number;

  const COL_MIN_WIDTH = 75;

  interface Col {
    id: string;
    name: string;
    el?: HTMLElement;
    min?: number;
  }

  let columns: Col[] = [
    {
      id: "branch",
      name: "Branch",
    },
    {
      id: "tree",
      name: "Tree",
      min: $tree.width,
    },
    {
      id: "subject",
      name: "Commit",
      min: COL_MIN_WIDTH,
    },
    {
      id: "gpg",
      name: "GPG",
      min: COL_MIN_WIDTH,
    },
    {
      id: "author",
      name: "Author",
      min: COL_MIN_WIDTH,
    },
    {
      id: "date",
      name: "Date",
      min: COL_MIN_WIDTH,
    },
  ];

  tree.subscribe(t => {
    columns[1].min = parseInt(t.width);
    if (columns[1].el) columns[1].el.style.minWidth = t.width;
  });

  onMount(() => {
    commitData.refresh();
    currentCommit.unset();

    commitData.subscribe(_ => {
      if (commitsTable) commitsTable.style.gridTemplateColumns = "auto auto auto auto auto auto";
    });
  });

  afterUpdate(() => {
    scrollDiv.scrollTo(0, scrollPosition);
  });

  function stickyScroll(el: HTMLElement) {
    scrollDiv = el;
    scrollDiv.addEventListener("scroll", () => scrollPosition = scrollDiv.scrollTop);
  }

  let w1: number;
  let w2: number;
  let startX: number;
  let max: number;
  let min: number;

  const resizeDown = (e: MouseEvent & { currentTarget: HTMLElement }) => {
    resizeIndex = parseInt(e.currentTarget.dataset.index);
    w1 = columns[resizeIndex].el.offsetWidth;
    w2 = columns[resizeIndex + 1].el.offsetWidth;
    startX = e.pageX;
    if (columns[resizeIndex].el.dataset.id === 'branch') {
      min = parseInt(window.getComputedStyle(columns[resizeIndex].el).width);
    } else {
      min = columns[resizeIndex].min;
    }
    max = w1 + w2 - columns[resizeIndex + 1].min;

    window.addEventListener('mousemove', resizeMove);
    window.addEventListener('mouseup', resizeUp);
  }

  const resizeMove = (e: MouseEvent) => {
    let gridColumns = columns.map(col => `${col.el.offsetWidth}px`);

    let widthLeft = Math.min(Math.max((w1 + e.pageX - startX), min), max);
    gridColumns[resizeIndex] = `${widthLeft}px`;
    if (resizeIndex + 1 < columns.length) {
      let widthRight = w2 + (w1 - widthLeft);
      gridColumns[resizeIndex + 1] = `${widthRight}px`;
    }

    commitsTable.style.gridTemplateColumns = `${gridColumns.join(' ')}`;
  }

  const resizeUp = () => {
    resizeIndex = null;
    window.removeEventListener('mousemove', resizeMove);
    window.removeEventListener('mouseup', resizeUp);
  }
</script>

<div class="commits" id="commits">
  <div class="commits__table-container" id="commits__scroll" use:stickyScroll use:setCommitsContainer>
    {#if Object.entries($commits).length}
      <table bind:this={commitsTable} class="commits__table" id="commits__table">
        <thead>
          <tr>
            {#each columns as col, i}
              <th class="commits__th commits__th--{col.id}"
                bind:this={columns[i].el}
                on:mousedown={resizeDown}
                data-index={i}
                data-id="{col.id}"
                style="min-width: {COL_MIN_WIDTH}px">
                {col.name}
                {#if i < columns.length - 1}
                  <div on:mousedown={resizeDown} data-index={i} class="resize-handle" style:height="24px"></div>
                {/if}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          <tr>
            <td></td>
            <td>
              <div class="tree">
                <div class="tree__graph" use:setFade data-fade="{$tree.continues}" data-height="{$tree.height}">{@html $tree.svg.outerHTML}</div>
              </div>
            </td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
          </tr>
          {#each Object.entries($commits) as [_, commit]}
            <CommitRow {commit} signStatus={$commitSignData.commits[commit.Hash]} />
          {/each}
        </tbody>
      </table>
      {#if $tree.continues}
        <LoadMoreCommits />
      {/if}
    {/if}
  </div>
  <CommitDetails />
</div>
