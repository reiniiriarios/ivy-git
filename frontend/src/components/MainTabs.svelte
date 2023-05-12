<script lang="ts">
  import Details from 'components/Details.svelte';
  import Diff from 'components/Diff.svelte';
  import Commits from 'components/Commits.svelte';

  let tab = "tree";
  (window as any).currentTab = tab;

  function tabChanges() {
    tab = "changes";
    (window as any).currentTab = tab;
  }

  function tabTree() {
    tab = "tree";
    (window as any).currentTab = tab;
  }

  function tabDetails() {
    tab = "details";
    (window as any).currentTab = tab;
  }
</script>

<nav class="tabs">
  <button class="tabs__tab {tab == 'changes' ? 'active' : ''}" id="tab-changes" on:click={tabChanges}>
    Changes
  </button>
  <button class="tabs__tab {tab == 'tree' ? 'active' : ''}" id="tab-tree" on:click={tabTree}>
    Commits
  </button>
  <button class="tabs__tab {tab == 'details' ? 'active' : ''}" id="tab-details" on:click={tabDetails}>
    Details
  </button>
</nav>

<Diff active={tab == 'changes'} />
<Commits active={tab == 'tree'} />
<Details active={tab == 'details'} />

<div class="window-resize window-resize--right"></div>
<div class="window-resize window-resize--bottom"></div>
