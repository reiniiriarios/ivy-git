<script lang="ts">
  import Info from "components/elements/Info.svelte";
  import Select from "components/elements/Select.svelte";
  import { settings } from "stores/settings";

  function saveTheme(e: CustomEvent) {
    settings.updateTheme(e.detail.value);
  }

  // todo: define these elsewhere
  let themes = {
    dark: 'Dark',
    light: 'Light',
  };
</script>

<div class="setting">
  <h4 class="setting__name">
    Theme
    <Info>
      More themes coming soon.
    </Info>
  </h4>
  <div class="setting__set">
    <Select
      options={themes}
      selected={$settings.Theme ?? ""}
      on:change={saveTheme}
    />
  </div>
</div>
