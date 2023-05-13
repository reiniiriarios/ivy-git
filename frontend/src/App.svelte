<script lang="ts">
  import "style/style.scss";
  import SelectRepo from "components/SelectRepo.svelte";
  import Confirm from 'components/Confirm.svelte';
  import Message from "components/Message.svelte";
  import SelectBranch from "components/SelectBranch.svelte";
  import Changes from "components/Changes.svelte";
  import MainTabs from "components/MainTabs.svelte";
  import TitleBar from "components/TitleBar.svelte";
  import ContextMenu from "components/ContextMenu.svelte";

  import { tabUpDown } from "scripts/keyboard-navigation";

  import { currentRepo, repos } from "stores/repos";
  import { branches, currentBranch } from "stores/branches";
  import { changes } from "stores/changes";
  import { currentTab } from "stores/current-tab";
  import { commitData } from "stores/commit-data";

  import type { EnvironmentInfo } from "wailsjs/runtime/runtime";
  import { ResizeWindow } from "wailsjs/go/main/App";

  // Load initial ui state.
  function init() {
    currentRepo.refresh();
    currentBranch.refresh();
    branches.refresh();
    repos.refresh();
    changes.refresh();
    if ($currentTab === 'tree') {
      commitData.refresh();
    }
    (window as any).runtime.Environment().then((env: EnvironmentInfo) => {
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

  window.addEventListener('keydown', tabUpDown);

  window.addEventListener('resize', ResizeWindow);
</script>

<TitleBar />
<div id="container">
  <div id="sidebar">
    <SelectRepo />
    <SelectBranch />
    <Changes />
  </div>
  <main>
    <MainTabs />
  </main>
  <Confirm />
  <Message />
</div>
<ContextMenu />
