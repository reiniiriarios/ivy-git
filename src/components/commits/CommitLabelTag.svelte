<script lang="ts">
  import octicons from "@primer/octicons";
  import type { Ref } from "stores/commit-data";

  export let tag: Ref;
</script>

<div class="refs__label refs__label--tag"
  title={tag.AbbrName == "" ? "" : tag.Name}
  data-name="{tag.Name}"
  data-menu="tag">
  <div class="refs__icon">{@html octicons['tag'].toSVG({ "width": 14 })}</div>
  <div class="refs__label-name">{tag.AbbrName != "" ? tag.AbbrName : tag.Name}</div>
  {#if tag.SyncedRemotes && tag.SyncedRemotes.length}
    {#each tag.SyncedRemotes as remote}
      <div class="refs__leaf">
        <span>{remote}</span>
      </div>
    {/each}
  {/if}
</div>
