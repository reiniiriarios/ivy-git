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

<div class="controls {os} {position}">
  {#if os == 'darwin'}
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
            <path d="M205 205 1843 205 1843 1843 205 1843 205 205ZM2 3-5 2048 2048 2048 2048-1 2-1Z"></path>
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
    <!-- These currently are the same as windows; update later. -->
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
            <path d="M205 205 1843 205 1843 1843 205 1843 205 205ZM2 3-5 2048 2048 2048 2048-1 2-1Z"></path>
          {/if}
        </svg>
      </span>
    </div>
      <div class="control" id="close" on:click={close} on:keyup={close}><span>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048">
          <path d="M1169 1024l879 879-145 145-879-879-879 879L0 1903l879-879L0 145 145 0l879 879L1903 0l145 145-879 879z"></path>
        </svg>
      </span>
    </div>
  {/if}
</div>
