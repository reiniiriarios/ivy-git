<script lang="ts">
  export let state: string = "";
</script>

<div class="getting-started">
  <div>
    {#if state === 'no-branch'}
      No branch selected. Select or create a branch.
    {:else if state === 'no-commits'}
      No commits have been added to this repo yet. Add a commit to enable the rest of the gui.
    {:else}
      Add or select a repo in the top left to get started.
    {/if}
  </div>
</div>
