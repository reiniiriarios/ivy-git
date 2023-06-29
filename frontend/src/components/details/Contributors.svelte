<script lang="ts">
  import octicons from "@primer/octicons";
  import Avatar from "components/elements/Avatar.svelte";
  import Info from "components/elements/Info.svelte";
  import { contributors } from "stores/contributors";
  import { currentRepo, repos } from "stores/repos";
  import { settings } from "stores/settings";
  import { currentTab } from "stores/ui";
  import { onMount } from "svelte";

  onMount(() => {
    contributors.fetch();
  });

  let updateButton: HTMLButtonElement;
  let goIcon: HTMLElement;
  let waitIcon: HTMLElement;

  function update() {
    updateButton.disabled = true;
    goIcon.classList.add('icon--hidden');
    waitIcon.classList.remove('icon--hidden');
    contributors.update().then(() => {
      // Artificial loading time here makes the UI make more sense.
      // This doesn't delay content loading.
      setTimeout(() => {
        updateButton.disabled = false;
        goIcon.classList.remove('icon--hidden');
        waitIcon.classList.add('icon--hidden');
      }, 200);
    });
  }
</script>

<div class="contributors">
  <div class="contributors__head">
    <div>
      <h2>
        Contributors
        <Info>
          Calculated by the number of commits on the main branch.
          Does not auto-refresh; can be cpu-intensive to calculate on larger repos.
        </Info>
      </h2>
    </div>
    {#if $repos[$currentRepo].Main}
      <div>
        <button
          class="btn btn--icon btn-sm"
          bind:this={updateButton}
          on:click={update}
        >
          Update
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
    {#if $contributors?.length}
      <ul class="contributors__list">
        {#each $contributors as c}
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
