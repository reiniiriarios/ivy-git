<script lang="ts">
  import octicons from '@primer/octicons';
  import { currentBranch } from 'stores/branches';
  import { type Ref } from 'stores/commits';

  export let branch: Ref;
  export let remotes: Ref[];
</script>

<div class="refs__label refs__label--branch"
  data-name="{branch.Name}"
  data-branch="{branch.Branch}"
  data-upstream="{branch.Upstream}"
  data-current="{$currentBranch.Name === branch.Name}"
  data-menu="branch">
  <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
  {#if branch.Head}
    <div class="refs__head">@</div>
  {/if}
  <div class="refs__label-name">{branch.Name}</div>
  {#if remotes?.length}
    {#each remotes as r}
      {#if r.Branch == branch.Branch}
        <div class="refs__leaf">
          {#if r.Head}
            <div class="refs__head">@</div>
          {/if}
          <span>{r.AbbrName != "" ? r.AbbrName : r.Remote}</span>
        </div>
      {/if}
    {/each}
  {/if}
</div>
