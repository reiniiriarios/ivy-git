<script lang="ts">
  import { createResizableColumn, setCommitsTable } from 'scripts/commit-table-resize';
  import CommitRow from 'components/CommitRow.svelte';
  import CommitDetails from 'components/CommitDetails.svelte';
  import { setCommitsContainer } from 'scripts/commit-details-resize';
  import { commits, tree } from 'stores/commit-data';
</script>

<div class="commits" id="commits">
  <div class="commits__table-container" use:setCommitsContainer>
    {#if Object.entries($commits).length}
      <table use:setCommitsTable class="commits__table" id="commits__table">
        <thead>
          <tr>
            <th use:createResizableColumn data-name="branch" data-order="0" class="commits__th commits__th--branch">Branch</th>
            <th use:createResizableColumn data-name="tree" data-order="1" class="commits__th commits__th--tree" style:min-width={$tree.width}>Tree</th>
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
                <div class="tree__graph">{@html $tree.svg.outerHTML}</div>
              </div>
            </td>
          </tr>
          {#each Object.entries($commits) as [_, commit]}
            <CommitRow {commit} />
          {/each}
        </tbody>
      </table>
    {/if}
  </div>
  <CommitDetails />
</div>
