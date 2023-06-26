<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">
    Highlight Main Branch
    <Info>
      Display a star next to the main branch in order to locate it easily.
    </Info>
  </h4>
  <Checkbox
    display="Highlight"
    checked={$settings.HighlightMainBranch}
    on:click={settings.toggleHighlightMainBranch}
  />
</div>
