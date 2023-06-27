<script lang="ts">
  import { onMount } from "svelte";

  import { gitConfig } from "stores/git-config";

  import SetWorkflow from "components/settings/SetWorkflow.svelte";
  import SetConventionalCommits from "components/settings/SetConventionalCommits.svelte";
  import SetMainBranch from "components/settings/SetMainBranch.svelte";
  import AppInfo from "components/settings/AppInfo.svelte";
  import SetGitConfig from "components/settings/SetGitConfig.svelte";
  import SetTheme from "components/settings/SetTheme.svelte";
  import SetHighlightMain from "components/settings/SetHighlightMain.svelte";
  import SetDisplayAvatars from "./SetDisplayAvatars.svelte";
  import SetDateFormat from "./SetDateFormat.svelte";

  onMount(() => {
    gitConfig.fetch();
  });
</script>

<div class="settings">
  <h2>Settings</h2>

  <h3>App</h3>
  <SetTheme />
  <SetWorkflow />
  <SetHighlightMain />
  <SetDateFormat />
  <SetConventionalCommits />
  <SetDisplayAvatars />

  <h3>Repo</h3>
  <SetMainBranch />

  <h3>Git</h3>
  <SetGitConfig />

  <h2>Info</h2>
  <AppInfo />
</div>
