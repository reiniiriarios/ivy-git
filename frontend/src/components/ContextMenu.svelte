<script lang="ts">
  import { menus, type Menu } from '../scripts/context-menus';

  const X_OFFSET = 3;

  let currentClickedElement: HTMLElement;

  function displayMenu(e: MouseEvent, menu: Menu) {
    let menuElement = document.getElementById("context-menu__" + menu.class);
    let x = e.pageX;
    let y = e.pageY;
    if (window.innerWidth - e.pageX < menuElement.offsetWidth + X_OFFSET) {
      x -= menuElement.offsetWidth - X_OFFSET;
    } else {
      x += X_OFFSET;
    }
    if (window.innerHeight - e.pageY < menuElement.offsetHeight) {
      y -= menuElement.offsetHeight;
    }
    menuElement.style.display = 'block';
    menuElement.style.left = x + "px";
    menuElement.style.top = y + "px";
  }

  function hideMenu(menu: Menu) {
    document.getElementById("context-menu__" + menu.class).style.display = 'none';
    currentClickedElement = null;
  }

  function hideMenus() {
    let menuElements = document.getElementsByClassName('context-menu');
    for (let i = 0; i < menuElements.length; i++) {
      (menuElements[i] as HTMLElement).style.display = 'none';
    }
    currentClickedElement = null;
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.addEventListener("contextmenu", function (e: MouseEvent & { target: HTMLElement }) {
      if (currentClickedElement) {
        e.preventDefault();
        hideMenus();
        return;
      }
      // todo: e.preventDefault() for everywhere in production mode
      for (let i = 0; i < menus.length; i++) {
        console.log(currentClickedElement);
        if (e.target.classList.contains(menus[i].class)) {
          currentClickedElement = e.target;
        } else {
          let n = e.target.parentNode;
          for (let j = 0; j < 4; j++, n = n.parentNode) {
            if (n instanceof HTMLElement && n.classList.contains(menus[i].class)) {
              currentClickedElement = n;
              break;
            }
          }
        }
        console.log(currentClickedElement);
        if (currentClickedElement) {
          e.preventDefault();
          console.log('displaying menu');
          displayMenu(e, menus[i]);
          break;
        } else {
          hideMenu(menus[i]);
        }
      };
    });

    document.body.addEventListener("keydown", function (e: KeyboardEvent) {
      if (e.key === "Escape") {
        hideMenus();
      }
    });

    document.addEventListener("click", function(e: MouseEvent) {
      hideMenus();
    });
  });
</script>

{#each menus as menu}
  <div class="context-menu" id="context-menu__{menu.class}">
    <ul class="context-menu__items">
      {#each menu.items as item}
        {#if item.text}
          <li class="context-menu__item">
            <div class="context-menu__action" on:click={(e) => item.callback(currentClickedElement)} on:keyup={() => {}}>{item.text}</div>
          </li>
        {:else if item.sep}
          <li class="context-menu__sep"></li>
        {/if}
      {/each}
    </ul>
  </div>
{/each}

<style lang="scss">
  .context-menu {
    display: none;
    position: absolute;
    z-index: 10000;
    min-width: 14rem;
    text-align: left;
    background-color: var(--color-context-bg);
    padding: 0.3rem;

    &__items {
      list-style: none;
      margin: 0;
      padding: 0;
    }

    &__item {
      display: block;
    }

    &__sep {
      margin: 0.4rem 0;
      height: 1px;
      background-color: var(--color-context-border);
    }

    &__action {
      cursor: pointer;
      padding: 0.5rem 1rem;

      &:hover {
        background-color: var(--color-context-bg-hover);
      }
    }
  }
</style>
