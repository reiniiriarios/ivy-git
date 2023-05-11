<script lang="ts">
  import { type FileStatDir } from "./CommitDetails.svelte";
  export let files: FileStatDir;
</script>

<span class="filestatdir__dir {!files?.Name ? 'filestatdir__dir--root' : files?.Files?.length || files?.Dirs?.length > 1 ? 'filestatdir__dir--files' : ''}">
  {#if files?.Name}
    {files.Name}
    {#if files?.Dirs?.length === 1 && !files?.Files?.length}
      /
    {/if}
  {/if}
  {#if files?.Dirs?.length}
    {#each Object.entries(files.Dirs) as [_, d]}
      <svelte:self files={d} />
    {/each}
  {/if}
  {#if files?.Files?.length}
    {#each Object.entries(files.Files) as [_, f]}
      <div class="filestatdir__file filestatdir__file--{f.Status}">
        {f.Name}
        <span class="diff">
          (<span class="added">+{f.Added}</span>, <span class="deleted">-{f.Deleted}</span>)
        </span>
      </div>
    {/each}
  {/if}
</span>
