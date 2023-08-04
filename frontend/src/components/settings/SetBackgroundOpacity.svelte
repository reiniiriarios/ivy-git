<script lang="ts">
  import { settings } from "stores/settings";

  function updateBackgroundOpacity(e: Event & { currentTarget: EventTarget & HTMLInputElement; }) {
    settings.updateBackgroundOpacity(parseInt(e.currentTarget.value));
  }
</script>

<div class="setting">
  <h4 class="setting__name">
    Background Opacity
  </h4>
  <div class="setting__set">
    <div class="range">
      <input
        type="range"
        class="range__slider"
        min="25" max="100" step="1"
        value={$settings.BackgroundOpacity}
        on:input={updateBackgroundOpacity}
      />
      <span class="range__value">{$settings.BackgroundOpacity}%</span>
    </div>
  </div>
</div>
