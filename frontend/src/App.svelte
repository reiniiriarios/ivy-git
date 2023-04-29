<script lang="ts">
  import "./style/style.scss";
  import SelectRepo from "./components/SelectRepo.svelte";
  import Confirm from './components/Confirm.svelte';
  import Message from "./components/Message.svelte";
  import SelectBranch from "./components/SelectBranch.svelte";
  import Changes from "./components/Changes.svelte";
  import MainTabs from "./components/MainTabs.svelte";
  import Tree from "./components/Tree.svelte";
  import Diff from "./components/Diff.svelte";
  import Details from "./components/Details.svelte";

  // Load initial ui state.
  function init() {
    (window as any).getSelectedRepo();
    (window as any).getRepos();
    (window as any).getCurrentBranch();
    (window as any).getBranches();
    (window as any).getChanges();
  }
  document.addEventListener('DOMContentLoaded', () => {
    init();
  });

  let tab = (window as any).currentTab ?? '';
</script>

<div id="sidebar">
  <SelectRepo />
  <SelectBranch />
  <Changes />
</div>
<main>
  <MainTabs />
  {#if tab == 'changes'}
    <Diff />
  {:else if tab == 'tree'}
    <Tree />
  {:else if tab == 'details'}
    <Details />
  {/if}
</main>
<Confirm />
<Message />
