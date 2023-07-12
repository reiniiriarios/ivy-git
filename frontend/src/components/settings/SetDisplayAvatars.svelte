<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import Info from "components/elements/Info.svelte";
  import { settings } from "stores/settings";
</script>

<div class="setting">
  <h4 class="setting__name">
    Avatars / Profile Photos
    <Info>
      Display avatars from Gravatar or GitHub for GitHub anon email addresses.
    </Info>
  </h4>
  <div class="setting__set">
    <Checkbox
      display="Display Avatars"
      checked={$settings.DisplayAvatars}
      on:click={settings.toggleDisplayAvatars}
    />
  </div>
</div>
