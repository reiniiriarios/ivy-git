<script lang="ts">
  import "./style/style.scss";
  import SelectRepo from "./components/SelectRepo.svelte";
  import Confirm from './components/Confirm.svelte';
  import Message from "./components/Message.svelte";
  import SelectBranch from "./components/SelectBranch.svelte";
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
