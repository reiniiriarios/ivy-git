<script lang="ts">
  import AppControl from 'components/app/AppControl.svelte';
  import type { EnvironmentInfo } from 'wailsjs/runtime/runtime';

  let os: string;
  (window as any).runtime.Environment().then((env: EnvironmentInfo) => {
    os = env.platform;
  });
</script>

<header class="titlebar">
  {#if !os}
    <div class="titlebar__title">Ivy Git</div>
  {:else if os == "darwin"}
    <AppControl os={os} position="left" />
    <div class="titlebar__title">Ivy Git</div>
  {:else}
    <div class="titlebar__title">Ivy Git</div>
    <AppControl os={os} position="right" />
  {/if}
</header>
