<script lang="ts">
  import octicons from "@primer/octicons";
  import { type Commit } from "stores/commits";

  export let commit: Commit;
</script>

<div class="refs__label refs__label--stash"
  data-hash="{commit.Hash}"
  data-ref="{commit.RefName}"
  data-subject="{window.btoa(commit.Subject)}"
  data-menu="stash">
  <div class="refs__icon">{@html octicons['inbox'].toSVG({ "width": 14 })}</div>
  <div class="refs__label-name">{commit.RefName}</div>
</div>
