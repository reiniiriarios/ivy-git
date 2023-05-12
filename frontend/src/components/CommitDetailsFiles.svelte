<script lang="ts">
  import { type FileStatDir } from "./CommitDetails.svelte";
  export let files: FileStatDir;
</script>

<div class="filestatdir__dir">
  {#if files?.Name}
    <div class="filestatdir__dir-name">
      {files.Name}
    </div>
  {/if}
  {#if files?.Dirs?.length}
    {#each Object.entries(files.Dirs) as [_, d]}
      <svelte:self files={d} />
    {/each}
  {/if}
  {#if files?.Files?.length}
    {#each Object.entries(files.Files) as [_, f]}
      <div class="filestatdir__file filestatdir__file--{f.Status}">
        {#if f.OldFile}
          <span class="filestatdir__file-old">
            {#if f.Dir === f.OldDir}
              {f.OldName}
            {:else}
              {f.OldFile}
            {/if}
            →
          </span>
        {/if}
        {f.Name}
        <span class="diff">
          (<span class="added">+{f.Added}</span>, <span class="deleted">-{f.Deleted}</span>)
        </span>
      </div>
    {/each}
  {/if}
</div>
