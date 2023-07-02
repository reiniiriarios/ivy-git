<script lang="ts">
  import Remotes from "components/details/Remotes.svelte";
  import RepoInfo from "components/details/RepoInfo.svelte";
  import CodeBreakdown from "components/details/CodeBreakdown.svelte";
  import Contributors from "components/details/Contributors.svelte";
</script>

<div class="details">
  <div class="details__inner">
    <Remotes />
    <RepoInfo />
    <CodeBreakdown />
    <Contributors />
  </div>
</div>
