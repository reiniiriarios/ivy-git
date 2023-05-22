<script lang="ts">
  import Details from 'components/Details.svelte';
  import Diff from 'components/Diff.svelte';
  import Commits from 'components/Commits.svelte';
  import Settings from 'components/Settings.svelte';
  import { currentTab } from 'stores/ui';
  import octicons from '@primer/octicons';
  import { currentFile } from 'stores/current-file';

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
    settings: {
      n: '<span aria-label="Settings">' + octicons.gear.toSVG({width: 18}) + '</span>',
      c: Settings,
    }
  };

  function setTab(t: string) {
    currentTab.set(t);
    currentFile.reset();
  }
</script>

<nav class="tabs">
  {#each Object.entries(tabs) as [t, d]}
    <button class="tabs__tab tabs__tab--{t}" class:active={$currentTab === t} id="tab-{t}" on:click={() => setTab(t)}>
      {@html d.n}
    </button>
  {/each}
</nav>

<svelte:component this={tabs[$currentTab].c} />

<div class="window-resize window-resize--right"></div>
<div class="window-resize window-resize--bottom"></div>
