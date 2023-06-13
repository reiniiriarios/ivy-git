<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">
    Conventional Commits
    <Info>
      See <a href="https://www.conventionalcommits.org/en/v1.0.0/" target="_blank">conventionalcommits.org</a> for more information.
    </Info>
  </h4>
  <Checkbox
    display="Highlight"
    checked={$settings.HighlightConventionalCommits}
    on:click={settings.toggleHighlightConventionalCommits}
  />
</div>
