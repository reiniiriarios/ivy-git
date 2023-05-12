<script lang="ts">
  import octicons from '@primer/octicons';

  export let os: string = "";
  export let position: string;

  let maximized = false;

  function minimize() {
    (window as any).runtime.WindowIsMinimised().then((isMinimized: boolean) => {
      isMinimized
        ? (window as any).runtime.WindowUnminimise()
        : (window as any).runtime.WindowMinimise();
    });
  }

  function maximize() {
    (window as any).runtime.WindowIsMaximised().then((isMaximized: boolean) => {
      isMaximized
        ? (window as any).runtime.WindowUnmaximise()
        : (window as any).runtime.WindowMaximise();
      maximized = !isMaximized;
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
  {:else if os == 'windows'}
    <div class="control" id="minimize" on:click={minimize} on:keyup={minimize}>
      <span>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048">
          <path d="M2048 819v205H0V819h2048z"></path>
        </svg>
      </span>
    </div>
    <div class="control" id="maximize" on:click={maximize} on:keyup={maximize}>
      <span>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048">
          {#if maximized}
            <path d="M2048 1638h-410v410H0V410h410V0h1638v1638zM1434 614H205v1229h1229V614zm409-409H614v205h1024v1024h205V205z"></path>
          {:else}
            <path d="M2048 0v819h-205V350L350 1843h469v205H0v-819h205v469L1698 205h-469V0h819z"></path>
          {/if}
        </svg>
      </span>
    </div>
    <div class="control" id="close" on:click={close} on:keyup={close}>
      <span>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048">
          <path d="M1169 1024l879 879-145 145-879-879-879 879L0 1903l879-879L0 145 145 0l879 879L1903 0l145 145-879 879z"></path>
        </svg>
      </span>
    </div>
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
      background-color: var(--color-scale-r-4-100);
    }

    #minimize {
      background-color: var(--color-scale-y-2-100);
    }

    #maximize {
      background-color: var(--color-scale-g-3-100);
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
        fill: var(--color-scale-a-4-100);

        svg {
          height: 10px;
          width: 10px;
        }
      }

      &:hover {
        background-color: var(--color-scale-a-6-100);
      }
    }

    #close:hover {
      background-color: var(--color-scale-r-5-100);

      span {
        fill: var(--color-scale-a-1-100);
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
        fill: var(--color-scale-a-4-100);
      }

      &:hover {
        span {
          fill: var(--color-scale-a-1-100);
        }
      }
    }
  }
</style>
