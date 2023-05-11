<script lang="ts">
  import Details from "./Details.svelte";
  import Diff from "./Diff.svelte";
  import Commits from "./Commits.svelte";

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
  <button class="tabs__tab {tab == 'changes' ? 'active' : ''}" id="tab-changes"
    on:click={tabChanges}
    on:keyup={tabChanges}>
    Changes
  </button>
  <button class="tabs__tab {tab == 'tree' ? 'active' : ''}" id="tab-tree"
    on:click={tabTree} on:keyup={tabTree}>
    Commits
  </button>
  <button
    class="tabs__tab {tab == 'details' ? 'active' : ''}" id="tab-details"
    on:click={tabDetails} on:keyup={tabDetails}>
    Details
  </button>
</nav>

<Diff active={tab == 'changes'} />
<Commits active={tab == 'tree'} />
<Details active={tab == 'details'} />

<div class="window-resize window-resize--right"></div>
<div class="window-resize window-resize--bottom"></div>
