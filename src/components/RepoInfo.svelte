<script lang="ts">
  import octicons from "@primer/octicons";
  import { numBranches, numTags, numCommits } from "stores/repo-info";
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
  </div>
</div>
