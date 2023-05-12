<script lang="ts">
  import octicons from '@primer/octicons';

  import { getLabelDist, type Commit, type Ref } from "../scripts/graph";
  import { currentBranch } from '../../src/stores/branches';

  function isCurrent(n: string): boolean {
    return $currentBranch.Name == n;
  }

  export let commit: Commit;
  export let HEAD: Ref;
</script>

<div class="refs">
  {#if commit.Branches && commit.Branches.length}
    {#each commit.Branches as b}
      <div class="refs__label refs__label--branch"
        data-name="{b.Name}"
        data-current="{isCurrent(b.Name)}"
        data-menu="branch">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{b.Name}</div>
        {#if commit.Remotes && commit.Remotes.length}
          {#each commit.Remotes as r}
            <div class="refs__leaf">{r.ShortName}</div>
          {/each}
        {/if}
      </div>
    {/each}
  {:else if commit.Remotes && commit.Remotes.length}
    {#each commit.Remotes as r}
      <div class="refs__label refs__label--branch"
        data-name="{r.Name}"
        data-remote="{r.ShortName}"
        data-menu="remoteBranch">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{r.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Hash == HEAD.Hash}
    <div class="refs__label refs__label--head"
      data-menu="head">
      <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
      <div class="refs__label-name">HEAD</div>
      {#if commit.Heads && commit.Heads.length}
        {#each commit.Heads as h}
          <div class="refs__leaf">{h.ShortName}</div>
        {/each}
      {/if}
    </div>
  {:else if commit.Heads && commit.Heads.length}
    {#each commit.Heads as h}
      <div class="refs__label refs__label--head"
        data-remote="{h.ShortName}"
        data-menu="remoteHead">
        <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{h.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Tags && commit.Tags.length}
    {#each commit.Tags as t}
      <div class="refs__label refs__label--tag"
        data-name="{t.Name}"
        data-menu="tag">
        <div class="refs__icon">{@html octicons['tag'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{t.Name}</div>
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

  <div class="refs__line" style="width:{getLabelDist(commit.X)}; right: -{getLabelDist(commit.X)};"></div>
</div>
