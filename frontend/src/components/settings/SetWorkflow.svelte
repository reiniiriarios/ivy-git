<script lang="ts">
  import octicons from "@primer/octicons";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";

  function saveWorkflow(workflow: string) {
    settings.updateWorkflow(workflow);
  }
</script>

<div class="setting">
  <h4 class="setting__name">
    Workflow
    <Info>
      This setting changes which items are available in some context menus.
    </Info>
  </h4>
  <div class="setting__set big-option">
    <button class="big-option__option"
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
    <button class="big-option__option"
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
    <button class="big-option__option"
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
