<script lang="ts">
  import octicons from "@primer/octicons";

  export let state: string = "";
</script>

<div class="getting-started">
  <div>
    {#if state === 'loading'}
      <div class="getting-started__loading">
        {@html octicons.gear.toSVG({width: 24})}
        <div>Loading app... Checking git...</div>
      </div>
    {:else if state === 'no-branch'}
      No branch selected. Select or create a branch.
    {:else if state === 'no-commits'}
      No commits have been added to this repo yet. Add a commit to enable the rest of the gui.
    {:else if state === 'no-git'}
      <p>Cannot find git executable. Please make sure Git is installed, then restart the app.</p>
      <p><a href="https://git-scm.com/downloads">Download Git</a></p>
    {:else}
      Add or select a repo in the top left to get started.
    {/if}
  </div>
</div>
