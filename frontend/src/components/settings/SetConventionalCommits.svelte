<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">Conventional Commits</h4>
  <Checkbox
    display="Highlight"
    checked={$settings.HighlightConventionalCommits}
    on:click={settings.toggleHighlightConventionalCommits}
  />
</div>
