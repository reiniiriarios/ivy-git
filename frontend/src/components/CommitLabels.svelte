<script lang="ts">
  import octicons from '@primer/octicons';

  import { getLabelDist } from 'scripts/graph';
  import { currentBranch } from 'stores/branches';
  import { HEAD, type Commit } from 'stores/commit-data';

  function isCurrent(n: string): boolean {
    return $currentBranch.Name == n;
  }

  export let commit: Commit;
</script>

<div class="refs">
  {#if commit.Branches && commit.Branches.length}
    {#each commit.Branches as b}
      <div class="refs__label refs__label--branch"
        data-name="{b.Name}"
        data-branch="{b.Branch}"
        data-upstream="{b.Upstream}"
        data-current="{isCurrent(b.Name)}"
        data-menu="branch">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{b.Name}</div>
        {#if commit.Remotes && commit.Remotes.length}
          {#each commit.Remotes as r}
            {#if r.Branch == b.Branch}
              <div class="refs__leaf">{r.AbbrName != "" ? r.AbbrName : r.Remote}</div>
            {/if}
          {/each}
        {/if}
      </div>
    {/each}
  {/if}
  {#if commit.Remotes && commit.Remotes.length}
    {#each commit.Remotes as r}
      <div class="refs__label refs__label--branch"
        title={r.AbbrName == "" ? "" : r.Name}
        data-name="{r.Name}"
        data-branch="{r.Branch}"
        data-remote="{r.Remote}"
        data-menu="remoteBranch">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{r.AbbrName != "" ? r.AbbrName : r.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Hash == $HEAD.Hash}
    <div class="refs__label refs__label--head"
      data-menu="head">
      <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
      <div class="refs__label-name">HEAD</div>
      {#if commit.Heads && commit.Heads.length}
        {#each commit.Heads as h}
          <div class="refs__leaf">{h.AbbrName != "" ? h.AbbrName : h.Remote}</div>
        {/each}
      {/if}
    </div>
  {:else if commit.Heads && commit.Heads.length}
    {#each commit.Heads as h}
      <div class="refs__label refs__label--head"
        title={h.AbbrName == "" ? "" : h.Name}
        data-remote="{h.Remote}"
        data-menu="remoteHead">
        <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{h.AbbrName != "" ? h.AbbrName : h.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Tags && commit.Tags.length}
    {#each commit.Tags as t}
      <div class="refs__label refs__label--tag"
        title={t.AbbrName == "" ? "" : t.Name}
        data-name="{t.Name}"
        data-menu="tag">
        <div class="refs__icon">{@html octicons['tag'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{t.AbbrName != "" ? t.AbbrName : t.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Stash}
    <div class="refs__label refs__label--stash"
      data-hash="{commit.Hash}"
      data-subject="{window.btoa(commit.Subject)}"
      data-menu="stash">
      <div class="refs__icon">{@html octicons['inbox'].toSVG({ "width": 14 })}</div>
      <div class="refs__label-name">{commit.RefName}</div>
    </div>
  {/if}

  <div class="refs__line" style:width={getLabelDist(commit.X)} style:right={'-'+getLabelDist(commit.X)}></div>
</div>
