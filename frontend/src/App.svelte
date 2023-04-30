<script lang="ts">
  import "./style/style.scss";
  import SelectRepo from "./components/SelectRepo.svelte";
  import Confirm from './components/Confirm.svelte';
  import Message from "./components/Message.svelte";
  import SelectBranch from "./components/SelectBranch.svelte";
  import Changes from "./components/Changes.svelte";
  import MainTabs from "./components/MainTabs.svelte";

  // Load initial ui state.
  function init() {
    (window as any).getSelectedRepo();
    (window as any).getRepos();
    (window as any).getCurrentBranch();
    (window as any).getBranches();
    (window as any).getChanges();
    if ((window as any).currentTab == 'tree') {
      (window as any).GetCommitsForTree();
    }
  }
  document.addEventListener('DOMContentLoaded', () => {
    init();
  });
</script>

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
