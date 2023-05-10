<script lang="ts">
  export let active: boolean;

  import { GetCommitList } from '../../wailsjs/go/main/App';
  import { drawGraph, getSVGWidth, type Commit, type Ref } from '../scripts/graph';
  import { createResizableColumn } from '../scripts/commit-table-resize';
  import CommitDetails from './CommitDetails.svelte';

  let commits: Commit[] = [];
  let HEAD: Ref;
  let svg: SVGSVGElement;
  let svgWidth: string;

  (window as any).GetCommitList = async () => {
    GetCommitList().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          commits = result.Commits as Commit[];
          HEAD = result.HEAD;
          console.log('HEAD', HEAD);
          console.log('commits', commits);
          svg = drawGraph(result.Graph);
          svgWidth = getSVGWidth(result.Graph);
          console.log('branches', result.Graph.Branches);
          console.log('vertices', result.Graph.Vertices);
          break;
      }
    });
  };

  const resize = (e: HTMLElement) => {
    createResizableColumn(e);
  }
</script>

{#if active}
  <div class="commits" id="commits">
    {#if Object.entries(commits).length}
      <table class="commits__table">
        <tr>
          <th use:resize class="commits__th commits__th--branch">Branch</th>
          <th use:resize class="commits__th commits__th--tree" style="min-width: {svgWidth};">
            <div class="commits__th-inner">Tree</div>
            <div class="tree">
              <div class="tree__graph">{@html svg.outerHTML}</div>
            </div>
          </th>
          <th use:resize class="commits__th commits__th--commit">Commit</th>
          <th use:resize class="commits__th commits__th--author">Author</th>
          <th class="commits__th commits__th--date">Date</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <CommitDetails commit={commit} HEAD={HEAD} />
        {/each}
      </table>
    {/if}
  </div>
{/if}

<style lang="scss">
  .commits {
    width: 100%;
    height: calc(100vh - var(--tabs-height) - var(--title-bar-height));
    overflow: auto;
    display: flex;
    flex-direction: row;
    justify-content: stretch;
    align-items: middle;

    &__table {
      margin-bottom: 0.5rem;
      width: 100%;
    }

    &__th {
      height: var(--commit-details-height);
      box-sizing: border-box;
      text-align: left;
      padding: 0.25rem 0.5rem 0.25rem 0.67rem;
      white-space: nowrap;
      background-color: var(--color-scale-a-7-100);
      cursor: default !important;
      user-select: none;
      -webkit-user-select: none;
      position: relative;

      &--branch {
        text-align: right;
        padding-right: 1rem;
      }

      &--tree {
        padding: 0;
      }

      &-inner {
        padding: 0.25rem 0.5rem;
      }

      &:first-child {
        border-right: 1px solid var(--color-scale-a-8-100);
      }

      &:not(:first-child) {
        border-left: 1px solid var(--color-scale-a-8-100);
      }
    }
  }
</style>
