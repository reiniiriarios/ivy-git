<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import Sidebar from './components/Sidebar.svelte';
  import Confirm from './components/Confirm.svelte';
</script>

<Confirm />
<Sidebar />
<main>
  ...
</main>

<style>
  main {
    z-index: 1;
    position: absolute;
    top: 0;
    left: var(--sidebar-width);
    width: calc(100% - var(--sidebar-width));
    height: 100%;
  }
</style>
