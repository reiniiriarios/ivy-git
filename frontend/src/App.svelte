<script lang="ts">
  import "./style/style.scss";
  import SelectRepo from "./components/SelectRepo.svelte";
  import Confirm from './components/Confirm.svelte';
  import Message from "./components/Message.svelte";
  import SelectBranch from "./components/SelectBranch.svelte";

  // Load initial ui state.
  function init() {
    (window as any).getSelectedRepo();
    (window as any).getRepos();
    (window as any).getCurrentBranch();
    (window as any).getBranches();
  }
  document.addEventListener('DOMContentLoaded', () => {
    init();
  });
</script>

<div id="sidebar">
  <SelectRepo />
  <SelectBranch />
</div>
<main>
  ...
</main>
<Confirm />
<Message />
