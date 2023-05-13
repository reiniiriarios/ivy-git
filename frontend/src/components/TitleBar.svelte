<script lang="ts">
  import AppControl from 'components/AppControl.svelte';
  import type { EnvironmentInfo } from 'wailsjs/runtime/runtime';

  let os: string;
  (window as any).runtime.Environment().then((env: EnvironmentInfo) => {
    os = env.platform;
  });
</script>

<header>
  {#if !os}
    <div id="title">Ivy Git</div>
  {:else if os == "darwin"}
    <AppControl os={os} position="left" />
    <div id="title">Ivy Git</div>
  {:else}
    <div id="title">Ivy Git</div>
    <AppControl os={os} position="right" />
  {/if}
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
