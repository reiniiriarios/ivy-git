<script lang="ts">
  import { toggleDir } from 'scripts/commit-details-collapse';
  import { type FileStatDir } from 'stores/commit-details';
  import { currentDiff } from 'stores/diffs';
  import { currentTab } from 'stores/ui';
  export let hash: string;
  export let files: FileStatDir;

  function fetchDiff(file: string, oldfile: string = "") {
    // Since we're switching tabs, clear the old diff away first.
    currentDiff.clear();
    currentDiff.fetchDiff(hash, file, oldfile);
    currentTab.set('changes');
  }
</script>

<div class="filestatdir__dir">
  {#if files?.Name}
    <div class="filestatdir__dir-name" on:click={toggleDir} on:keypress={toggleDir}>
      {files.Name}
    </div>
  {/if}
  {#if files?.Dirs?.length}
    {#each Object.entries(files.Dirs) as [_, d]}
      <!-- RECURSION HERE -->
      <svelte:self hash={hash} files={d} />
    {/each}
  {/if}
  {#if files?.Files?.length}
    {#each Object.entries(files.Files) as [_, f]}
      <div
        class="filestatdir__file filestatdir__file--{f.Status}"
        on:click={() => fetchDiff(f.File, f.OldFile)}
        on:keypress={() => fetchDiff(f.File, f.OldFile)}
      >
        {#if f.OldFile}
          <span class="filestatdir__file-old">
            {#if f.Dir === f.OldDir}
              {f.OldName}
            {:else}
              {#if f.OldRel}{f.OldRel}/{/if}{f.OldName}
            {/if}
            →
          </span>
        {/if}
        {f.Name}
        {#if !f.Binary && f.Added || f.Deleted}
          <span class="filestatdir__diff">
            ({#if f.Status !== 'D'}<span class="added">+{f.Added}</span>{/if}{#if !['A','D'].includes(f.Status)}, {/if}{#if f.Status !== 'A'}<span class="deleted">-{f.Deleted}</span>{/if})
          </span>
        {/if}
      </div>
    {/each}
  {/if}
</div>
