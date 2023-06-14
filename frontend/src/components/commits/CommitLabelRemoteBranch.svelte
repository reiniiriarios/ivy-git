<script lang="ts">
  import octicons from "@primer/octicons";
  import type { Ref } from "stores/commits";

  export let remote: Ref;
</script>

<div class="refs__label refs__label--branch"
  title={remote.AbbrName == "" ? "" : remote.Name}
  data-name="{remote.Name}"
  data-branch="{remote.Branch}"
  data-remote="{remote.Remote}"
  data-menu="remoteBranch">
  <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
  <div class="refs__leaf">
    {#if remote.Head}
      <div class="refs__head">@</div>
    {/if}
    <span>{remote.AbbrName != "" ? remote.AbbrName : remote.Name}</span>
  </div>
</div>
