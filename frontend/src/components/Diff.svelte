<script lang="ts">
  export let active: boolean;
</script>

{#if active}
  <div class="diff">

  </div>
{/if}

<style lang="scss">
  .diff {
    width: 100%;
    height: 100%;
  }
</style>
