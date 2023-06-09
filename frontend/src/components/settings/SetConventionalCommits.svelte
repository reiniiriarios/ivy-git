<script lang="ts">
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">Conventional Commits</h4>
  <label class="checkbox">
    <input type="checkbox" checked={$settings.HighlightConventionalCommits} on:click={settings.toggleHighlightConventionalCommits}>
    <span></span>
    Highlight
  </label>
</div>
