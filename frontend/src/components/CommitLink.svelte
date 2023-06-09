<script lang="ts">
  import { commitsMap, commits } from "stores/commit-data";
  import { currentCommit } from "stores/commit-details";
  import { currentTab } from "stores/ui";
  import { get } from "svelte/store";
  import { onMount } from "svelte";

  export let hash: string;
  let commit_id: number;
  let linkable: boolean;

  $: commit_id = get(commitsMap).get(hash);
  $: linkable = typeof commit_id === 'number';

  function setCurrentCommit() {
    if (linkable && $commits[commit_id]) {
      currentTab.set('tree');
      currentCommit.set($commits[commit_id]);
    }
  }
</script>

<span
  class="linked-commit"
  class:linked-commit--linkable={linkable}
  on:click={setCurrentCommit}
  on:keypress={setCurrentCommit}
><slot>{hash}</slot></span>
