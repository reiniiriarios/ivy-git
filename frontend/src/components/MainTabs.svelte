<script lang="ts">
  import Details from 'components/Details.svelte';
  import Diff from 'components/Diff.svelte';
  import Commits from 'components/Commits.svelte';
  import { currentTab } from 'stores/current-tab';

  let tabs = {
    changes: {
      n:'Changes',
      c: Diff,
    },
    tree: {
      n:'Commits',
      c: Commits,
    },
    details: {
      n:'Details',
      c: Details,
    },
  };
</script>

<nav class="tabs">
  {#each Object.entries(tabs) as [t, d]}
    <button class="tabs__tab {$currentTab == t ? 'active' : ''}" id="tab-{t}" on:click={() => $currentTab = t}>
      {d.n}
    </button>
  {/each}
</nav>

<svelte:component this={tabs[$currentTab].c} />

<div class="window-resize window-resize--right"></div>
<div class="window-resize window-resize--bottom"></div>
