<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">
    Conventional Commits
    <Info>
      Display conventional commit prefix tags in different colors in the commits list.
      See <a href="https://www.conventionalcommits.org/en/v1.0.0/" target="_blank">conventionalcommits.org</a> for more information.
    </Info>
  </h4>
  <div class="setting__set">
    <Checkbox
      display="Highlight Commit Prefixes"
      checked={$settings.HighlightConventionalCommits}
      on:click={settings.toggleHighlightConventionalCommits}
    />
  </div>
</div>
