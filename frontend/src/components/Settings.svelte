<script lang="ts">
  import octicons from "@primer/octicons";
  import { branches } from "stores/branches";
  import { currentRepo, repos } from "stores/repos";
  import { settings } from "stores/settings";
  import Select from "./Select.svelte";

  function saveWorkflow(workflow: string) {
    settings.updateWorkflow(workflow);
  }

  function saveMain(e: CustomEvent) {
    repos.updateMain(e.detail.value);
  }
</script>

<div class="settings">
  <h2>Settings</h2>

  <h3>App</h3>

  <div class="setting">
    <h4 class="setting__name">Workflow</h4>
    <div class="big-option">
      <button class="btn big-option__option"
        class:big-option__option--active={$settings.Workflow === 'merge'}
        on:click={() => saveWorkflow('merge')}>
        <div class="big-option__icon">
          {@html octicons["git-merge"].toSVG({width: 24})}
        </div>
        <div class="big-option__name">
          Merge Commits
        </div>
        <div class="big-option__desc">
          Avoid editing history.
        </div>
      </button>
      <button class="btn big-option__option"
        class:big-option__option--active={$settings.Workflow === 'squash'}
        on:click={() => saveWorkflow('squash')}>
        <div class="big-option__icon">
          {@html octicons["git-merge"].toSVG({width: 24})}
        </div>
        <div class="big-option__name">
          Squash & Rebase
        </div>
        <div class="big-option__desc">
          Merge or rebase to update, squash to merge.
        </div>
      </button>
      <button class="btn big-option__option"
        class:big-option__option--active={$settings.Workflow === 'rebase'}
        on:click={() => saveWorkflow('rebase')}>
        <div class="big-option__icon">
          {@html octicons["git-merge"].toSVG({width: 24})}
        </div>
        <div class="big-option__name">
          Rebase Everything
        </div>
        <div class="big-option__desc">
          An advanced workflow for a simplified timeline.
        </div>
      </button>
    </div>
  </div>

  <div class="setting">
    <h4 class="setting__name">Conventional Commits</h4>
    <label class="checkbox">
      <input type="checkbox" checked={$settings.HighlightConventionalCommits} on:click={settings.toggleHighlightConventionalCommits}>
      <span></span>
      Highlight
    </label>
  </div>

  <h3>Repo</h3>

  <div class="setting">
    <h4 class="setting__name">Main Branch</h4>
    <Select values={$branches.map(b => b.Name)} selected={$repos[$currentRepo].Main ?? ""} on:change={saveMain} />
  </div>

  <h2>Info</h2>

  <div class="setting">
    <h4 class="setting__name">Version</h4>
    <div>{$settings.Version}</div>
  </div>

  <div class="setting">
    <h4 class="setting__name">Author</h4>
    <div><a href="https://github.com/reiniiriarios">Emma Litwa-Vulcu</a></div>
  </div>

  <div class="setting">
    <h4 class="setting__name">Help</h4>
    <div>
      <a href="https://github.com/reiniiriarios/ivy-git">Homepage</a> (todo)
      <span class="text-sep"></span>
      <a href="https://github.com/reiniiriarios/ivy-git">Report Issue</a> (todo)
    </div>
  </div>
</div>
