<script lang="ts">
  import { afterUpdate, onDestroy, onMount } from 'svelte';

  import CommitRow from 'components/commits/CommitRow.svelte';
  import CommitDetails from 'components/commit-details/CommitDetails.svelte';
  import LoadMoreCommits from 'components/commits/LoadMoreCommits.svelte';

  import { setCommitsContainer } from 'scripts/commit-details-resize';

  import { commitData, commits, tree, commitSignData } from 'stores/commits';
  import { setFade } from 'scripts/graph';
  import { settings } from 'stores/settings';
  import type { Unsubscriber } from 'svelte/store';

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
      min: parseInt($tree.width),
    },
    {
      id: "subject",
      name: "Commit",
      min: COL_MIN_WIDTH,
    },
  ];
  if ($settings.DisplayCommitSignatureInList) {
    columns.push({
      id: "gpg",
      name: "GPG",
      min: COL_MIN_WIDTH,
    });
  }
  columns = columns.concat([
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
  ]);

  const treeUnsubscribe = tree.subscribe(t => {
    columns[1].min = parseInt(t.width);
    if (columns[1].el) columns[1].el.style.minWidth = t.width;
  });

  let commitDataUnsubscribe: Unsubscriber;

  onMount(() => {
    commitData.refresh();
    commitDataUnsubscribe = commitData.subscribe(() => setAutoCols());
  });

  onDestroy(() => {
    treeUnsubscribe();
    commitDataUnsubscribe();
  });

  afterUpdate(() => {
    scrollDiv.scrollTo(0, scrollPosition);
  });

  function setAutoCols() {
    if (commitsTable) {
      commitsTable.style.gridTemplateColumns =
        $settings.DisplayCommitSignatureInList
          ? "auto auto 5fr auto auto auto"
          : "auto auto 5fr auto auto";
    }
  }

  function stickyScroll(el: HTMLElement) {
    scrollDiv = el;
    scrollDiv.addEventListener("scroll", () => scrollPosition = scrollDiv.scrollTop);
  }

  let w1: number;
  let w2: number;
  let startX: number;
  let max: number;
  let min: number;
  let resizer: HTMLElement;

  const resizeDown = (e: MouseEvent & { currentTarget: HTMLElement }) => {
    resizer = e.currentTarget;
    resizer.classList.add('active');

    resizeIndex = parseInt(resizer.dataset.index);
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
    resizer.classList.remove('active');
    resizeIndex = null;
    resizer = null;
    window.removeEventListener('mousemove', resizeMove);
    window.removeEventListener('mouseup', resizeUp);
  }
</script>

<div class="commits" id="commits">
  <div class="commits__table-container" id="commits__scroll" use:stickyScroll use:setCommitsContainer>
    {#if Object.entries($commits).length}
      <table bind:this={commitsTable} class="commits__table" id="commits__table"
        style="grid-template-columns: {$settings.DisplayCommitSignatureInList ? "auto auto 5fr auto auto auto" : "auto auto 5fr auto auto"}">
        <thead>
          <tr>
            {#each columns as col, i}
              {#if col.id !== 'gpg' || $settings.DisplayCommitSignatureInList}
                <th class="commits__th commits__th--{col.id}"
                  bind:this={columns[i].el}
                  data-index={i}
                  data-id="{col.id}"
                  style="min-width: {col.min ?? 0}px">
                  {col.name}
                  {#if i < columns.length - 1}
                    <div on:mousedown={resizeDown} data-index={i} class="resize-handle" style:height="24px"></div>
                  {/if}
                </th>
              {/if}
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
            {#if $settings.DisplayCommitSignatureInList}<td></td>{/if}
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
