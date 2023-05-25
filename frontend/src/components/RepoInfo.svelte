<script lang="ts">
  import octicons from "@primer/octicons";
  import { numBranches, numTags, numCommits, cloc } from "stores/repo-info";
  import { currentRepo, repos } from "stores/repos";
  import { onMount } from "svelte";
  import languages from "style/languages.json"

  onMount(() => {
    numBranches.fetch();
    numTags.fetch();
    numCommits.fetch();
    cloc.fetch();
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
  <h2>Code</h2>
  <div class="code-breakdown">
    {#if $cloc.Languages?.length}
      <div class="code-breakdown__bar">
        {#each Object.entries($cloc.Languages) as [_, lang]}
          <div
            style:background-color={languages[lang.Name] ?? '#ccc'}
            style:width={lang.TotalPercent.toFixed(4) + '%'}
          ></div>
        {/each}
      </div>
      <table>
        <thead>
          <tr>
            <th></th>
            <th aria-label="Language"></th>
            <th>Total Lines</th>
            <th>Code Lines</th>
            <th>Comments</th>
            <th>Blank Lines</th>
            <th>Files</th>
            <th>Percent</th>
          </tr>
        </thead>
        <tbody>
          {#each Object.entries($cloc.Languages) as [_, lang]}
            <tr>
              <td class="dot"><span style:background-color={languages[lang.Name] ?? '#ccc'}></span></td>
              <td class="name">{lang.Name}</td>
              <td>{lang.Total}</td>
              <td>{lang.Code}</td>
              <td>{lang.Comments}</td>
              <td>{lang.Blanks}</td>
              <td>{lang.Files}</td>
              <td>{lang.TotalPercent.toFixed(2)}%</td>
            </tr>
          {/each}
        </tbody>
      </table>
    {:else if $cloc.Error}
      <div class="error">{$cloc.Error}</div>
    {:else}
      <div class="loading">...</div>
    {/if}
  </div>
</div>
