<script lang="ts">
  import octicons from "@primer/octicons";
  import { numBranches, numTags, numCommits, cloc } from "stores/repo-info";
  import { currentRepo, repos } from "stores/repos";
  import { onMount } from "svelte";

  onMount(() => {
    numBranches.fetch();
    numTags.fetch();
    numCommits.fetch();
  })
</script>

<div class="repo-info">
  <h2>Info</h2>
  <div class="repo-info__things">
    <div>
      {@html octicons["git-branch"].toSVG({width: 16})}
      <strong>{$numBranches}</strong>
      {$numBranches === 1 ? 'branch' : 'branches'}
    </div>
    <div>
      {@html octicons["tag"].toSVG({width: 16})}
      <strong>{$numTags}</strong>
      {$numTags === 1 ? 'tag' : 'tags'}
    </div>
    <div>
      {@html octicons["git-commit"].toSVG({width: 16})}
      <strong>{$numCommits}</strong>
      {$numCommits === 1 ? 'commit' : 'commits'}
      on
      <strong>{$repos[$currentRepo].Main}</strong>
    </div>
    {#if $cloc.Total?.Files}
      <div>
        {@html octicons["file"].toSVG({width: 16})}
        <strong>{$cloc.Total.Files}</strong>
        {$cloc.Total.Files === 1 ? 'file' : 'files'}
      </div>
    {/if}
    {#if $cloc.Total?.Total}
      <div>
        {@html octicons["code"].toSVG({width: 16})}
        <strong>{$cloc.Total.Total}</strong>
        {$cloc.Total.Total === 1 ? 'line' : 'lines'}
        of code
      </div>
    {/if}
  </div>
</div>
