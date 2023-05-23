<script lang="ts">
  import "style/style.scss";
  import SelectRepo from "components/SelectRepo.svelte";
  import Message from "components/Message.svelte";
  import SelectBranch from "components/SelectBranch.svelte";
  import Changes from "components/Changes.svelte";
  import MainTabs from "components/MainTabs.svelte";
  import TitleBar from "components/TitleBar.svelte";
  import ContextMenu from "components/ContextMenu.svelte";
  import GetStarted from "components/GetStarted.svelte";
  import MakeCommit from "components/MakeCommit.svelte";
  import RemoteActions from "components/RemoteActions.svelte";

  import { addInputListener, keyboardNavListener } from "scripts/keyboard-navigation";
  import { addLinkListener } from "scripts/links";
  import { envInit, getPlatform } from "scripts/env";

  import { currentRepo, repos } from "stores/repos";
  import { currentRemote, remoteData } from "stores/remotes";
  import { settings } from "stores/settings";

  import { ResizeWindow } from "wailsjs/go/main/App";
  import { enableWatcher } from "events/watcher";

  // Load initial ui state.
  function init() {
    currentRepo.refresh();
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

  window.addEventListener('resize', ResizeWindow);
  keyboardNavListener();
  addLinkListener();
  addInputListener();

  enableWatcher();

  // Fixes an issue on macOS where when dragging the cursor will change to
  // the text selector. By only attaching this to HTMLElements, text itself
  // is still selectable.
  window.addEventListener('selectstart', (e: Event) => {
    if (e.target instanceof HTMLElement) {
      e.preventDefault();
    }
  });

  // Development: If hot updating a module, re-init the app for correct data cascade.
  if (import.meta.hot) {
    import.meta.hot.on('vite:afterUpdate', () => {
      window.location.reload();
    });
  }
</script>

<TitleBar />
<div id="container">
  <div class="sidebar">
    <SelectRepo />
    {#if $currentRepo}
      <SelectBranch />
    {/if}
    <Changes />
    {#if $currentRepo}
      {#if $currentRemote}
        <RemoteActions />
      {/if}
      <MakeCommit />
    {/if}
  </div>
  <main>
    {#if $currentRepo}
      <MainTabs />
    {:else}
      <GetStarted />
    {/if}
  </main>
  <Message />
</div>
<ContextMenu />
