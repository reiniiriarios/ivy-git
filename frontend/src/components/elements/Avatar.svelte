<script lang="ts">
  import { avatars } from "stores/avatars";

  let url: string = '';
  export let email: string = '';
  $: avatars.fetch(email).then(u => url = u);

  export let hover: boolean = false;
</script>

{#if url}
  <span class="avatar">
    <img class="avatar__small" src="{url}" alt="" />
    {#if hover}
      <div class="avatar__hover-box">
        <img class="avatar__big" src="{url}" alt="" />
      </div>
    {/if}
  </span>
{/if}
