<script lang="ts">
  import { avatars } from "stores/avatars";

  let url: string = '';
  export let email: string = '';
  $: avatars.fetch(email).then(u => url = u);
</script>

{#if url}
  <span class="avatar">
    <img src="{url}" alt="" />
  </span>
{/if}
