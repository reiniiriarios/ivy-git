<script lang="ts">
  import { noBranchSelected } from "stores/branches";
  import { currentRemote } from "stores/remotes";
  import { currentRepo } from "stores/repos";

  import Changes from "components/sidebar/Changes.svelte";
  import MakeCommit from "components/sidebar/MakeCommit.svelte";
  import RemoteActions from "components/sidebar/RemoteActions.svelte";
  import RepoStateBanner from "components/sidebar/RepoStateBanner.svelte";
  import SelectBranch from "components/sidebar/SelectBranch.svelte";
  import SelectRepo from "components/sidebar/SelectRepo.svelte";


</script>

<div class="sidebar">
  <SelectRepo />
  {#if $currentRepo}
    <SelectBranch />
  {/if}
  <RepoStateBanner />
  <Changes />
  {#if $currentRepo}
    {#if $currentRemote?.Name}
      <RemoteActions />
    {/if}
    {#if !$noBranchSelected}
      <MakeCommit />
    {/if}
  {/if}
</div>
