<script lang="ts">
  import { noBranchSelected } from "stores/branches";
  import { currentRemote } from "stores/remotes";
  import { currentRepo } from "stores/repos";
  import { RepoState, repoState } from "stores/repo-state";

  import Changes from "components/sidebar/Changes.svelte";
  import MakeCommit from "components/sidebar/MakeCommit.svelte";
  import RemoteActions from "components/sidebar/RemoteActions.svelte";
  import RepoStateBanner from "components/sidebar/RepoStateBanner.svelte";
  import SelectBranch from "components/sidebar/SelectBranch.svelte";
  import SelectRepo from "components/sidebar/SelectRepo.svelte";
  import RebaseActions from "components/sidebar/RebaseActions.svelte";
  import CherryPickActions from "components/sidebar/CherryPickActions.svelte";

  import { resizableSidebar } from "scripts/sidebar-resize";
</script>

<div class="sidebar" use:resizableSidebar>
  <SelectRepo />
  {#if $currentRepo}
    <SelectBranch />
  {/if}
  <RepoStateBanner />
  <Changes />
  {#if $currentRepo}
    {#if [RepoState.RebaseMerge, RepoState.ApplyOrRebase, RepoState.Interactive].includes($repoState)}
      <RebaseActions />
    {:else if [RepoState.CherryPick, RepoState.CherryPickSequence].includes($repoState)}
      <CherryPickActions />
    {:else if $currentRemote?.Name}
      <RemoteActions />
    {/if}
    {#if !$noBranchSelected}
      <MakeCommit />
    {/if}
  {/if}
</div>
