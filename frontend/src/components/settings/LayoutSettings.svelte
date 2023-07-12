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
  <h2 class="settings__title">Settings</h2>

    <h3 class="settings__section-header">App</h3>
    <div class="settings__grid">
      <SetTheme />
    </div>
    <SetWorkflow />
    <div class="settings__grid">
      <SetHighlightMain />
      <SetDateFormat />
      <SetConventionalCommits />
      <SetDisplayAvatars />
    </div>

  <div class="settings__section">
    <h3 class="settings__section-header">Repo</h3>
    <div class="settings__grid">
      <SetMainBranch />
    </div>
  </div>

  <div class="settings__section">
    <h3 class="settings__section-header">Git</h3>
    <SetGitConfig />
  </div>

  <div class="settings__section settings__section--info">
    <h2 class="settings__section-header">Info</h2>
    <div class="settings__grid">
      <AppInfo />
    </div>
  </div>
</div>
