<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">
    Auto Fetch
    <Info>
      Enables auto-fetch for current repo and current remote.
    </Info>
  </h4>
  <div class="setting__set">
    <Checkbox
      display="Enable Auto Fetch"
      checked={$settings.AutoFetch}
      on:click={settings.toggleAutoFetch}
    />
  </div>
</div>
