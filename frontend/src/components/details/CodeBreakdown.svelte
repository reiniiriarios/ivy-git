<script lang="ts">
  import octicons from "@primer/octicons";
  import { cloc } from "stores/cloc";
  import { onMount } from "svelte";
  import languages from "style/languages.json"
  import { formatBytes } from "scripts/bytes";
  import Info from "components/elements/Info.svelte";
  import { currentRepo, repos } from "stores/repos";
  import { currentTab } from "stores/ui";

  onMount(() => {
    cloc.fetch();
  });

  let updateButton: HTMLButtonElement;
  let goIcon: HTMLElement;
  let waitIcon: HTMLElement;

  let updateWord: string = "Update";
  let updateMessages: string[] = [
    "Working",
    "Still Working",
    "Large Repo, Please Wait",
    "Just a bit longer...",
    "Please hold...",
    "Wow, this is... a lot...",
    "*Sigh*",
  ];

  let commitsBehind: number = 0;
  cloc.subscribe(() => {
    cloc.numCommitsBehind().then(n => commitsBehind = n);
  });

  let running: boolean = false;

  function update() {
    running = true;
    updateButton.disabled = true;
    goIcon.classList.add('icon--hidden');
    waitIcon.classList.remove('icon--hidden');
    let done = false;
    let messageCounter = 0;
    cloc.update().then(() => {
      // Artificial loading time here makes the UI make more sense.
      // This doesn't delay content loading.
      setTimeout(() => {
        updateButton.disabled = false;
        goIcon.classList.remove('icon--hidden');
        waitIcon.classList.add('icon--hidden');
        done = true;
        updateWord = 'Update';
        running = false;
      }, 200);
    });
    // Set an update timer to give update messages to make the
    // user feel okay about waiting a long time. Because this can
    // take a Long Time on Big Repos, particularly on slower machines.
    let update = setInterval(() => {
      if (done) {
        clearInterval(update);
        return;
      }
      if (messageCounter === updateMessages.length) {
        messageCounter = 0;
      }
      updateWord = updateMessages[messageCounter];
      messageCounter++;
    }, 3000);
  }
</script>

<div class="code-breakdown">
  <div class="code-breakdown__head">
    <div>
      <h2>
        Code
        <Info>
          May not be exact.
          Some auto-generated files, such as <code>package-lock.json</code>, are ignored.
          Calculated based on main branch.
          Does not auto-refresh; can be cpu-intensive to calculate on larger repos.
        </Info>
      </h2>
    </div>
    {#if commitsBehind}
      <div class="contributors__behind">
        {commitsBehind} {commitsBehind === 1 ? 'commit' : 'commits'} out of date.
      </div>
    {/if}
    {#if $repos[$currentRepo].Main}
      <div>
        <button
          class="btn btn--icon btn-sm"
          bind:this={updateButton}
          on:click={update}
        >
          {updateWord}
          <span class="icon" bind:this={goIcon}>
            {@html octicons["arrow-switch"].toSVG({width: 12})}
          </span>
          <span class="icon icon--hidden icon--spin" bind:this={waitIcon}>
            {@html octicons["gear"].toSVG({width: 12})}
          </span>
        </button>
      </div>
    {/if}
  </div>
  <div class="code-breakdown__inner">
    {#if $repos[$currentRepo].Main}
      {#if $cloc.Languages && Object.entries($cloc.Languages).length}
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
              <th>Size</th>
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
                <td>{formatBytes(lang.Bytes)}</td>
                <td>{lang.TotalPercent.toFixed(2)}%</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {:else if $cloc.Error}
        <div class="code-breakdown__message code-breakdown__message--error">{$cloc.Error}</div>
      {:else if running}
        <div class="code-breakdown__message code-breakdown__message--loading">{@html octicons.gear.toSVG({width: 24})}</div>
      {:else}
        <div class="code-breakdown__message">
          No data.
          <div>
            <button class="btn" on:click={update}>Update Contributors Data</button>
          </div>
        </div>
      {/if}
    {:else}
      <div class="code-breakdown__message">
        Select a main branch in in order to view code breakdown.
        <div>
          <button class="btn" on:click={() => currentTab.set('settings')}>View Settings</button>
        </div>
      </div>
    {/if}
  </div>
</div>
