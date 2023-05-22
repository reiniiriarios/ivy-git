<script lang="ts">
  import octicons from "@primer/octicons";
  import { numBranches, numTags, mainInfo } from "stores/repo-info";
  import { onMount } from "svelte";

  onMount(() => {
    numBranches.fetch();
    numTags.fetch();
    mainInfo.fetch();
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
      <strong>{$mainInfo.Count}</strong>
      {$mainInfo.Count === 1 ? 'commit' : 'commits'}
      on
      <strong>{$mainInfo.Name}</strong>
    </div>
  </div>
</div>
