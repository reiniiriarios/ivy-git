<script lang="ts">
  import octicons from "@primer/octicons";

  export let os: string = "";
  export let position: string;

  function minimize() {
    (window as any).runtime.WindowIsMinimised().then((isMinimized: Boolean) => {
      isMinimized
        ? (window as any).runtime.WindowUnminimise()
        : (window as any).runtime.WindowMinimise();
    });
  }

  function maximize() {
    (window as any).runtime.WindowIsMaximised().then((isMaximized: Boolean) => {
      isMaximized
        ? (window as any).runtime.WindowUnmaximise()
        : (window as any).runtime.WindowMaximise();
    });
  }

  function close() {
    (window as any).runtime.Quit();
  }
</script>

<div id="controls" class="{os} {position}">
  {#if position == 'left'}
    <div class="control" id="close" on:click={close} on:keyup={close}><span>{@html octicons.x.toSVG({ "width": 13 })}</span></div>
    <div class="control" id="minimize" on:click={minimize} on:keyup={minimize}><span>{@html octicons.dash.toSVG({ "width": 13 })}</span></div>
    <div class="control" id="maximize" on:click={maximize} on:keyup={maximize}><span>{@html octicons.plus.toSVG({ "width": 13 })}</span></div>
  {:else}
    <div class="control" id="minimize" on:click={minimize} on:keyup={minimize}><span>{@html octicons.dash.toSVG({ "width": 14 })}</span></div>
    <div class="control" id="maximize" on:click={maximize} on:keyup={maximize}><span>{@html octicons['screen-full'].toSVG({ "width": 13 })}</span></div>
    <div class="control" id="close" on:click={close} on:keyup={close}><span>{@html octicons.x.toSVG({ "width": 16 })}</span></div>
  {/if}
</div>

<style lang="scss">
  #controls {
    position: absolute;
    top: 0;
    left: 0;
    display: flex;

    &.right {
      right: 0;
      left: initial;
    }
  }

  .control {
    --wails-draggable: no-drag;
  }

  .darwin {
    padding: 0.25rem 0.5rem;
    box-sizing: border-box;

    .control {
      width: 1rem;
      height: 1rem;
      border-radius: 1rem;
      margin: 0.25rem 0.3rem;
      box-sizing: border-box;

      span {
        display: none;
        fill: #000;
        opacity: 0.5;
        mix-blend-mode: multiply;
      }
    }

    #close {
      background-color: var(--color-scale-red-4);
    }

    #minimize {
      background-color: var(--color-scale-yellow-2);
    }

    #maximize {
      background-color: var(--color-scale-green-3);
    }

    &:hover {
      .control span {
        display: block;
      }
    }
  }

  .windows {
    .control {
      width: calc(var(--title-bar-height) * 1.5);
      height: var(--title-bar-height);

      span {
        display: flex;
        height: var(--title-bar-height);
        justify-content: center;
        align-items: center;
        fill: var(--color-scale-gray-4);
      }

      &:hover {
        background-color: var(--color-scale-gray-6);
      }
    }

    #close:hover {
      background-color: var(--color-scale-red-5);

      span {
        fill: var(--color-scale-gray-1);
      }
    }
  }

  .linux {
    .control {
      width: calc(var(--title-bar-height) * 1.25);
      height: var(--title-bar-height);

      span {
        display: flex;
        height: var(--title-bar-height);
        justify-content: center;
        align-items: center;
        fill: var(--color-scale-gray-4);
      }

      &:hover {
        span {
          fill: var(--color-scale-gray-1);
        }
      }
    }
  }
</style>
