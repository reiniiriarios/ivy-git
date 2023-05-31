<script lang="ts">
  import octicons from "@primer/octicons";
  import { cloc } from "stores/repo-info";
  import { onMount } from "svelte";
  import languages from "style/languages.json"

  onMount(() => {
    cloc.fetch();
  })
</script>

<div class="code-breakdown">
  <h2>Code</h2>
  <div class="code-breakdown__inner">
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
      <div class="code-breakdown__error">{$cloc.Error}</div>
    {:else}
      <div class="code-breakdown__loading">{@html octicons.gear.toSVG({width: 24})}</div>
    {/if}
  </div>
</div>