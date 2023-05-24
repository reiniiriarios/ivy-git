<script lang="ts">
  import AppControl from 'components/AppControl.svelte';

  let os: string = 'darwin';
  // (window as any).runtime.Environment().then((env: EnvironmentInfo) => {
  //   os = env.platform;
  // });
</script>

<header data-tauri-drag-region>
  {#if !os}
    <div data-tauri-drag-region id="title">Ivy Git</div>
  {:else if os == "darwin"}
    <AppControl os={os} position="left" />
    <div data-tauri-drag-region id="title">Ivy Git</div>
  {:else}
    <div data-tauri-drag-region id="title">Ivy Git</div>
    <AppControl os={os} position="right" />
  {/if}
</header>

<style lang="scss">
  @use "../style/!mixins" as *;

  header {
    --wails-draggable: drag;
    @include no-select;

    height: var(--title-bar-height);
    width: 100vw;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    padding-bottom: 0.1rem;

    #title {
      @include no-select;
      width: 8rem;
    }
  }
</style>
