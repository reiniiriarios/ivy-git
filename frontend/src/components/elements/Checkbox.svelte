<script lang="ts">
  export let checked: boolean = false;
  export let display: string;
  export let id: string = "";
</script>

<label class="checkbox">
  <input type="checkbox" {id} checked={checked} on:click>
  <span></span>
  {display}
</label>
