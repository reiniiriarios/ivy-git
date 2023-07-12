<script lang="ts">
  import Select from "components/elements/Select.svelte";
  import { dateFormats, settings } from "stores/settings";

  function saveDateFormat(e: CustomEvent) {
    settings.updateDateFormat(parseInt(e.detail.value));
  }

  let opts: {[val: number]: string} = {};
  $: Object.entries($dateFormats).forEach(([i, v]) => {
    opts[i] = v.Display;
  });
</script>

<div class="setting">
  <h4 class="setting__name">Date Format</h4>
  <div class="setting__set">
    <Select
      options={opts}
      selected={$settings.DateFormat ?? 0}
      on:change={saveDateFormat}
    />
  </div>
</div>
