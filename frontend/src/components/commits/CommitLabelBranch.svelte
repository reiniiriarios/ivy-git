<script lang="ts">
  import octicons from '@primer/octicons';
  import { currentBranch } from 'stores/branches';
  import { type Ref } from 'stores/commits';
  import { currentRepo, repos } from 'stores/repos';
  import { settings } from 'stores/settings';

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
  {#if $settings.HighlightMainBranch && branch.Name === $repos[$currentRepo].Main}
    <div class="refs__main" aria-label="Main branch">
      {@html octicons['star-fill'].toSVG({width: 14})}
    </div>
  {/if}
  {#if branch.Head}
    <div class="refs__head" aria-label="Head">@</div>
  {/if}
  <div class="refs__label-name">{branch.Name}</div>
  {#if remotes?.length}
    {#each remotes as r}
      {#if r.Branch == branch.Branch}
        <div class="refs__leaf">
          {#if r.Head}
            <div class="refs__head" aria-label="Remote Head">@</div>
          {/if}
          <span>{r.AbbrName != "" ? r.AbbrName : r.Remote}</span>
        </div>
      {/if}
    {/each}
  {/if}
</div>
