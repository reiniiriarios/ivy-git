<script lang="ts">
  import LayoutDetails from 'components/details/LayoutDetails.svelte';
  import Commits from 'components/commits/Commits.svelte';
  import LayoutSettings from 'components/settings/LayoutSettings.svelte';
  import { currentTab } from 'stores/ui';
  import octicons from '@primer/octicons';
  import LayoutDiffs from 'components/diffs/LayoutDiffs.svelte';

  let tabs = {
    changes: {
      n:'Changes',
      c: LayoutDiffs,
    },
    tree: {
      n:'Commits',
      c: Commits,
    },
    details: {
      n:'Details',
      c: LayoutDetails,
    },
    settings: {
      n: '<span class="tabs__icon" aria-label="Settings">' + octicons.gear.toSVG({width: 18}) + '</span>',
      c: LayoutSettings,
    }
  };

  function setTab(t: string) {
    currentTab.set(t);
  }
</script>

<nav class="tabs">
  {#each Object.entries(tabs) as [t, d]}
    <button class="tabs__tab tabs__tab--{t}" class:active={$currentTab === t} id="tab-{t}" on:click={() => setTab(t)}>
      {@html d.n}
    </button>
  {/each}
</nav>
<div class="tab-content">
  <svelte:component this={tabs[$currentTab].c} />
</div>

<div class="window-resize window-resize--right"></div>
<div class="window-resize window-resize--bottom"></div>
