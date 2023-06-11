<script lang="ts">
  import SetWorkflow from "components/settings/SetWorkflow.svelte";
  import SetConventionalCommits from "components/settings/SetConventionalCommits.svelte";
  import SetMainBranch from "components/settings/SetMainBranch.svelte";
  import Info from "components/settings/Info.svelte";
  import SetGitConfig from "./SetGitConfig.svelte";
  import { onMount } from "svelte";
  import { gitConfig } from "stores/git-config";

  onMount(() => {
    gitConfig.fetch();
  });
</script>

<div class="settings">
  <h2>Settings</h2>

  <h3>App</h3>
  <SetWorkflow />
  <SetConventionalCommits />

  <h3>Repo</h3>
  <SetMainBranch />

  <h3>Git</h3>
  <SetGitConfig />

  <h2>Info</h2>
  <Info />
</div>
