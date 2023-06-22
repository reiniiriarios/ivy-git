<script lang="ts">
  import "style/style.scss";
  import Message from "components/messages/Message.svelte";
  import MainTabs from "components/app/MainTabs.svelte";
  import TitleBar from "components/app/TitleBar.svelte";
  import ContextMenu from "components/ContextMenu.svelte";
  import GetStarted from "components/GetStarted.svelte";

  import { GitIsInstalled, ResizeWindow } from "wailsjs/go/main/App";

  import { addInputListener, keyboardNavListener } from "scripts/keyboard-navigation";
  import { addLinkListener } from "scripts/links";
  import { envInit } from "scripts/env";

  import { appData } from "stores/app-data";
  import { currentRepo, repos } from "stores/repos";
  import { remoteData } from "stores/remotes";
  import { settings, theme } from "stores/settings";
  import { noBranchSelected } from "stores/branches";

  import { enableWatcher } from "events/watcher";
  import LayoutSidebar from "components/sidebar/LayoutSidebar.svelte";

  let gitInstalled: boolean = false;
  GitIsInstalled().then(r => gitInstalled = r);

  let goos: string = "";

  // Load initial ui state.
  appData.fetch();
  currentRepo.load();
  repos.refresh();
  settings.refresh();
  remoteData.refresh();
  envInit().then(env => {
    goos = env.platform;
    switch (env.platform) {
      case "darwin":
        document.documentElement.style.setProperty("--color-app-bg", "var(--color-app-bg--darwin)");
        break;
      case "windows":
        document.documentElement.style.setProperty("--color-app-bg", "var(--color-app-bg--windows)");
        break;
    }
  });

  // Keep theme updated in <html data-theme="theme">
  theme.subscribe(t => {
    document.documentElement.dataset.theme = t;
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

<div class="app app--{goos}">
  <TitleBar />
  <div id="container">
    {#if !gitInstalled}
      <GetStarted state="no-git" />
    {:else}
      <LayoutSidebar />
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
    {/if}
    <Message />
  </div>
  <ContextMenu />
</div>
