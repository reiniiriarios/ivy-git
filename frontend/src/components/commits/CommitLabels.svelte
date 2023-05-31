<script lang="ts">
  import { getLabelDist } from 'scripts/graph';
  import { HEAD, type Commit } from 'stores/commit-data';
  import CommitLabelBranch from './CommitLabelBranch.svelte';
  import CommitLabelRemoteBranch from './CommitLabelRemoteBranch.svelte';
  import CommitLabelTag from './CommitLabelTag.svelte';
  import CommitLabelStash from './CommitLabelStash.svelte';
  import CommitLabelHeads from './CommitLabelHeads.svelte';

  export let commit: Commit;
</script>

<div class="refs">
  {#if commit.Branches && commit.Branches.length}
    {#each commit.Branches as b}
      <CommitLabelBranch branch={b} remotes={commit.RemoteBranches} />
    {/each}
  {/if}
  {#if commit.RemoteBranches && commit.RemoteBranches.length}
    {#each commit.RemoteBranches as r}
      {#if !r.SyncedLocally}
        <CommitLabelRemoteBranch remote={r} />
      {/if}
    {/each}
  {/if}

  {#if commit.Tags && commit.Tags.length}
    {#each commit.Tags as t}
      <CommitLabelTag tag={t} />
    {/each}
  {/if}

  {#if commit.Stash}
    <CommitLabelStash commit={commit} />
  {/if}

  {#if (!commit.Branches?.length && commit.Hash === $HEAD.Hash) || (!commit.RemoteBranches?.length && commit.Heads?.length)}
    <CommitLabelHeads isHEAD={commit.Hash === $HEAD.Hash} heads={commit.Heads} />
  {/if}

  <div class="refs__line" style:width={getLabelDist(commit.X)} style:right={'-'+getLabelDist(commit.X)}></div>
</div>
