<script lang="ts">
  import SetWorkflow from "components/settings/SetWorkflow.svelte";
  import SetConventionalCommits from "components/settings/SetConventionalCommits.svelte";
  import SetMainBranch from "components/settings/SetMainBranch.svelte";
  import Info from "components/settings/Info.svelte";
</script>

<div class="settings">
  <h2>Settings</h2>

  <h3>App</h3>
  <SetWorkflow />
  <SetConventionalCommits />

  <h3>Repo</h3>
  <SetMainBranch />

  <h2>Info</h2>
  <Info />
</div>
