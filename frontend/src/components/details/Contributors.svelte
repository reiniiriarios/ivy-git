<script lang="ts">
  import octicons from "@primer/octicons";
  import Avatar from "components/elements/Avatar.svelte";
  import Info from "components/elements/Info.svelte";
  import { contributors } from "stores/contributors";
  import { currentRepo, repos } from "stores/repos";
  import { settings } from "stores/settings";
  import { currentTab } from "stores/ui";
  import { onDestroy, onMount } from "svelte";

  onMount(() => {
    contributors.fetch();
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
  const contributorsUnsubscribe = contributors.subscribe(() => {
    contributors.numCommitsBehind().then(n => commitsBehind = n);
  });

  onDestroy(() => {
    contributorsUnsubscribe();
  });

  let running: boolean = false;

  function update() {
    running = true;
    updateButton.disabled = true;
    goIcon.classList.add('icon--hidden');
    waitIcon.classList.remove('icon--hidden');
    let done = false;
    let messageCounter = 0;
    contributors.update().then(() => {
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

<div class="contributors">
  <div class="contributors__head" data-menu="contributors">
    <div>
      <h2>
        Contributors
        <Info>
          Calculated by the number of commits on the main branch.
          Does not auto-refresh; can be cpu-intensive to calculate on larger repos.
          Respects <a href="https://git-scm.com/docs/gitmailmap">mailmap</a>.
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
  {#if $repos[$currentRepo].Main}
    {#if $contributors?.Contributors?.length}
      <ul class="contributors__list">
        {#each $contributors.Contributors as c}
          <li class="contributors__contributor">
            {#if $settings.DisplayAvatars && c.Email}
              <div class="contributors__avatar">
                <Avatar email="{c.Email}" />
              </div>
            {/if}
            <div>
              <div class="contributors__name">{c.Name}</div>
              <div class="contributors__email">
                <a href="mailto:{c.Email}">{c.Email}</a>
              </div>
              <div class="contributors__details">
                <span class="contributors__commits">{c.Commits.toLocaleString("en-US")} commits</span>
                <span class="contributors__insertions">{c.Insertions.toLocaleString("en-US")} ++</span>
                <span class="contributors__deletions">{c.Deletions.toLocaleString("en-US")} --</span>
              </div>
            </div>
          </li>
        {/each}
      </ul>
    {:else if !running}
      <div class="contributors__message">
        The contributors list displays users' contributions to main.<br>
        Calculating on larger repos can take time.
        <div>
          <button class="btn" on:click={update}>Update Contributors Data</button>
        </div>
      </div>
    {/if}
  {:else}
    <div class="contributors__message">
      Select a main branch in in order to view contributors.
      <div>
        <button class="btn" on:click={() => currentTab.set('settings')}>View Settings</button>
      </div>
    </div>
  {/if}
</div>
