<script lang="ts">
  import Details from "./Details.svelte";
  import Diff from "./Diff.svelte";
  import Tree from "./Tree.svelte";

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
    Tree
  </button>
  <button
    class="tabs__tab {tab == 'details' ? 'active' : ''}" id="tab-details"
    on:click={tabDetails} on:keyup={tabDetails}>
    Details
  </button>
</nav>

<Diff active={tab == 'changes'} />
<Tree active={tab == 'tree'} />
<Details active={tab == 'details'} />

<style lang="scss">
  .tabs {
    width: 100%;
    display: flex;
    border-bottom: 1px solid var(--color-btn-border);
    height: var(--tabs-height);

    &__tab {
      width: 33.333%;
      border: none;
      color: var(--color-text);
      background-color: var(--color-btn-bg);
      font-size: 1rem;
      padding: 0.5rem 1rem;
      margin: 0;
      cursor: pointer;

      &:not(:last-child) {
        border-right: 1px solid var(--color-btn-border);
      }

      &.active {
        background-color: var(--color-btn-bg-active);
      }

      &:hover {
        background-color: var(--color-btn-bg-hover);
      }
    }
  }
</style>
