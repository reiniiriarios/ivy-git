<script lang="ts">
  import "style/style.scss";
  import SelectRepo from "components/sidebar/SelectRepo.svelte";
  import Message from "components/messages/Message.svelte";
  import SelectBranch from "components/sidebar/SelectBranch.svelte";
  import Changes from "components/sidebar/Changes.svelte";
  import MainTabs from "components/app/MainTabs.svelte";
  import TitleBar from "components/app/TitleBar.svelte";
  import ContextMenu from "components/ContextMenu.svelte";
  import GetStarted from "components/GetStarted.svelte";
  import MakeCommit from "components/sidebar/MakeCommit.svelte";
  import RemoteActions from "components/sidebar/RemoteActions.svelte";
  import RepoStateBanner from "components/sidebar/RepoStateBanner.svelte";

  import { ResizeWindow } from "wailsjs/go/main/App";

  import { addInputListener, keyboardNavListener } from "scripts/keyboard-navigation";
  import { addLinkListener } from "scripts/links";
  import { envInit } from "scripts/env";

  import { appData } from "stores/app-data";
  import { currentRepo, repos } from "stores/repos";
  import { currentRemote, remoteData } from "stores/remotes";
  import { settings } from "stores/settings";
  import { noBranchSelected } from "stores/branches";

  import { enableWatcher } from "events/watcher";

  // Load initial ui state.
  function init() {
    appData.fetch();
    currentRepo.load();
    repos.refresh();
    settings.refresh();
    remoteData.refresh();
    envInit().then(env => {
      switch (env.platform) {
        case "darwin":
          document.documentElement.style.setProperty("--color-app-bg", "var(--color-app-bg--darwin)");
          break;
        case "windows":
          document.documentElement.style.setProperty("--color-app-bg", "var(--color-app-bg--windows)");
          break;
      }
    });
  }
  document.addEventListener('DOMContentLoaded', () => {
    init();
  });

  // Frontend Listeners
  window.addEventListener('resize', ResizeWindow);
  keyboardNavListener();
  addLinkListener();
  addInputListener();

  // Fixes an issue on macOS where when dragging the cursor will change to
  // the text selector. By only attaching this to HTMLElements, text itself
  // is still selectable.
  window.addEventListener('selectstart', (e: Event) => {
    if (e.target instanceof HTMLElement) {
      e.preventDefault();
    }
  });

  // Backend Listeners
  enableWatcher();
</script>

<TitleBar />
<div id="container">
  <div class="sidebar">
    <SelectRepo />
    {#if $currentRepo}
      <SelectBranch />
    {/if}
    <RepoStateBanner />
    <Changes />
    {#if $currentRepo}
      {#if $currentRemote?.Name}
        <RemoteActions />
      {/if}
      {#if !$noBranchSelected}
        <MakeCommit />
      {/if}
    {/if}
  </div>
  <main>
    {#if $currentRepo}
      {#if $noBranchSelected}
        <GetStarted state="no-branch" />
      {:else}
        <MainTabs />
      {/if}
    {:else}
      <GetStarted />
    {/if}
  </main>
  <Message />
</div>
<ContextMenu />
