<script lang="ts">
  import CodeBreakdown from "components/details/CodeBreakdown.svelte";
  import Remotes from "components/details/Remotes.svelte";
  import RepoInfo from "components/details/RepoInfo.svelte";
</script>

<div class="details">
  <div class="details__inner">
    <Remotes />
    <RepoInfo />
    <CodeBreakdown />
  </div>
</div>
