<script lang="ts">
  import { branches } from "stores/branches";
  import { currentRepo, repos } from "stores/repos";
  import Select from "components/Select.svelte";

  function saveMain(e: CustomEvent) {
    repos.updateMain(e.detail.value);
  }
</script>

<div class="setting">
  <h4 class="setting__name">Main Branch</h4>
  <Select
    values={$branches.map(b => b.Name)}
    selected={$repos[$currentRepo].Main ?? ""}
    on:change={saveMain}
  />
</div>
