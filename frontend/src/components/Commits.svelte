<script lang="ts">
  export let active: boolean;

  import { GetCommitList } from 'wailsjs/go/main/App';
  import { drawGraph, getSVGWidth, type Commit, type Ref } from 'scripts/graph';
  import { createResizableColumn, setCommitsTable } from 'scripts/commit-table-resize';
  import CommitRow from 'components/CommitRow.svelte';
  import CommitDetails from 'components/CommitDetails.svelte';
  import { setCommitsContainer } from 'scripts/commit-details-resize';

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
</script>

{#if active}
  <div class="commits" id="commits">
    <div class="commits__table-container" use:setCommitsContainer>
      {#if Object.entries(commits).length}
        <table use:setCommitsTable class="commits__table" id="commits__table">
          <thead>
            <tr>
              <th use:createResizableColumn data-name="branch" data-order="0" class="commits__th commits__th--branch">Branch</th>
              <th use:createResizableColumn data-name="tree" data-order="1" class="commits__th commits__th--tree" style="min-width: {svgWidth};">Tree</th>
              <th use:createResizableColumn data-name="subject" data-order="2" class="commits__th commits__th--subject">Commit</th>
              <th use:createResizableColumn data-name="authorName" data-order="3" class="commits__th commits__th--author">Author</th>
              <th use:createResizableColumn data-name="authorDate" data-order="4" data-resizeflex class="commits__th commits__th--date">Date</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td></td>
              <td>
                <div class="tree">
                  <div class="tree__graph">{@html svg.outerHTML}</div>
                </div>
              </td>
            </tr>
            {#each Object.entries(commits) as [_, commit]}
              <CommitRow commit={commit} HEAD={HEAD} />
            {/each}
          </tbody>
        </table>
      {/if}
    </div>
    <CommitDetails />
  </div>
{/if}
