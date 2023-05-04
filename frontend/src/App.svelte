<script lang="ts">
  import "./style/style.scss";
  import SelectRepo from "./components/SelectRepo.svelte";
  import Confirm from './components/Confirm.svelte';
  import Message from "./components/Message.svelte";
  import SelectBranch from "./components/SelectBranch.svelte";
  import Changes from "./components/Changes.svelte";
  import MainTabs from "./components/MainTabs.svelte";
  import TitleBar from "./components/TitleBar.svelte";
  import { GoOs } from "../wailsjs/go/main/App";

  // Load initial ui state.
  function init() {
    (window as any).getSelectedRepo();
    (window as any).getRepos();
    (window as any).getCurrentBranch();
    (window as any).getBranches();
    (window as any).getChanges();
    if ((window as any).currentTab == 'tree') {
      (window as any).GetCommitList();
    }
    GoOs().then(os => {
      switch (os) {
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
