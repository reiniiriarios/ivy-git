<script lang="ts">
  export let active: boolean;

  import { GetCommitList } from '../../wailsjs/go/main/App';

  import { drawGraph, UNCOMMITED_HASH } from '../scripts/graph';
  import type { Commit, Ref } from '../scripts/graph';
  import CommitDetails from './CommitDetails.svelte';
  import CommitLabels from './CommitLabels.svelte';

  let commits: Commit[] = [];
  let HEAD: Ref;
  let svg: SVGSVGElement;

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
          console.log('branches', result.Graph.Branches);
          console.log('vertices', result.Graph.Vertices);
          break;
      }
    });
  };
</script>

{#if active}
  <div class="commits" id="commits">
    {#if Object.entries(commits).length}
      <table class="commits__branches">
        <tr>
          <th class="h-b">Branch</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <tr class="c-{commit.Color} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''}">
            <td>
              {#if commit.Labeled}
                <CommitLabels commit={commit} HEAD={HEAD} />
              {/if}
            </td>
          </tr>
        {/each}
      </table>
      <div id="tree" class="tree">
        <div class="tree__text">Tree</div>
        <div class="tree__graph">{@html svg.outerHTML}</div>
      </div>
      <table class="commits__details">
        <tr>
          <th class="h-c">Commit</th>
          <th>Author</th>
          <th>Date</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <CommitDetails commit={commit} />
        {/each}
      </table>
    {/if}
  </div>
{/if}

<style lang="scss">
  .commits {
    min-width: 100%;
    height: calc(100vh - var(--tabs-height) - var(--title-bar-height));
    overflow: auto;
    display: flex;
    flex-direction: row;
    justify-content: stretch;
    align-items: middle;

    table {
      margin-bottom: 0.5rem;

      th, td {
        height: var(--commit-details-height);
        box-sizing: border-box;
      }

      th {
        text-align: left;
        padding: 0.25rem 0.5rem;
        white-space: nowrap;
        background-color: var(--color-scale-a-7-100);

        &.h-b {
          text-align: right;
          padding-right: 1rem;
        }

        &.h-c {
          width: 100%;
        }
      }
    }

    &__branches {
      th {
        border-right: 1px solid var(--color-scale-a-8-100);
      }

      td {
        padding-left: 0.67rem;
      }
    }

    &__details {
      flex: 1;

      th {
        border-left: 1px solid var(--color-scale-a-8-100);
      }
    }
  }
</style>
