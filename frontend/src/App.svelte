<script lang="ts">
  import "style/style.scss";
  import Message from "components/messages/Message.svelte";
  import MainTabs from "components/app/MainTabs.svelte";
  import TitleBar from "components/app/TitleBar.svelte";
  import ContextMenu from "components/ContextMenu.svelte";
  import GetStarted from "components/GetStarted.svelte";
  import LayoutSidebar from "components/sidebar/LayoutSidebar.svelte";
  import LinkPreview from "components/elements/LinkPreview.svelte";

  import { GitIsInstalled, ResizeWindow } from "wailsjs/go/ivy/App";

  import { addInputListener, keyboardNavListener } from "scripts/keyboard-navigation";
  import { addLinkListener } from "scripts/links";
  import { setMainBlock } from "scripts/sidebar-resize";

  import { get } from "svelte/store";

  import { environment } from "stores/env";
  import { appData } from "stores/app-data";
  import { currentRepo, repos } from "stores/repos";
  import { remoteData } from "stores/remotes";
  import { settings, theme } from "stores/settings";
  import { noBranchSelected } from "stores/branches";

  import { enableWatcher } from "events/watcher";
  import { autoFetchTimer } from "events/auto-fetch";
  import { onDestroy } from "svelte";

  let checkingApp: boolean = true;
  let gitInstalled: boolean = false;
  GitIsInstalled().then(r => {
    gitInstalled = r;
    checkingApp = false;
  });

  let goos: string = "";

  // Load initial ui state.
  appData.fetch();
  currentRepo.load();
  repos.refresh();
  settings.refresh().then(() => {
    document.documentElement.style.setProperty("--bg-opacity", (get(settings).BackgroundOpacity ?? 100) + '%');
    if (get(settings).AutoFetch) {
      autoFetchTimer.init();
    }
  });
  remoteData.refresh();
  environment.fetch();

  // Keep theme updated in <html data-theme="theme">
  const themeUnsubscribe = theme.subscribe(t => {
    document.documentElement.dataset.theme = t;
  });
  onDestroy(() => {
    themeUnsubscribe();
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
    {#if checkingApp}
      <GetStarted state="loading" />
    {:else if !gitInstalled}
      <GetStarted state="no-git" />
    {:else}
      <LayoutSidebar />
      <main use:setMainBlock>
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
  <LinkPreview />
  <ContextMenu />
</div>
