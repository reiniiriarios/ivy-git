<script lang="ts">
  import type { Ref } from "stores/commits";

  export let isHEAD: boolean;
  export let heads: Ref[];
</script>

{#if isHEAD}
  <div class="refs__label refs__label--head"
    data-menu="head">
    <div class="refs__icon">@</div>
    <div class="refs__label-name">HEAD</div>
    {#if heads?.length}
      {#each heads as h}
        <div class="refs__leaf"><span>{h.AbbrName != "" ? h.AbbrName : h.Remote}</span></div>
      {/each}
    {/if}
  </div>
{:else}
  {#each heads as h}
    <div class="refs__label refs__label--head"
      title={h.AbbrName == "" ? "" : h.Name}
      data-remote="{h.Remote}"
      data-menu="remoteHead">
      <div class="refs__icon">@</div>
      <div class="refs__leaf"><span>{h.AbbrName != "" ? h.AbbrName : h.Name}</span></div>
    </div>
  {/each}
{/if}
