<script lang="ts">
  import { GoOs } from "../../wailsjs/go/main/App";
  import AppControl from "./AppControl.svelte";
</script>

<header>
  {#await GoOs()}
    <div id="title">Ivy Git</div>
  {:then os}
    {#if os == "darwin"}
      <AppControl os={os} position="left" />
      <div id="title">Ivy Git</div>
    {:else}
      <div id="title">Ivy Git</div>
      <AppControl os={os} position="right" />
    {/if}
  {:catch}
    <AppControl position="right" />
  {/await}
</header>

<style lang="scss">
  header {
    --wails-draggable: drag;
    cursor: default !important;
    user-select: none;
    -webkit-user-select: none;

    height: var(--title-bar-height);
    width: 100vw;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    padding-bottom: 0.1rem;

    #title {
      user-select: none;
      width: 8rem;
    }
  }
</style>
