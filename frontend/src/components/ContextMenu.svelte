<script lang="ts">
  import { menus, type MenuItem } from 'context-menus/_all';

  const X_OFFSET = 3;

  let currentClickedElement: HTMLElement;
  let currentMenu: string;
  let menuItems: MenuItem[] = [];

  function displayMenu(e: MouseEvent) {
    let menuElement = document.getElementById("context-menu");
    let x = e.pageX;
    let y = e.pageY;
    let w = getCurrentMenuWidth();
    let h = getCurrentMenuHeight();
    if (window.innerWidth - e.pageX < w + X_OFFSET) {
      x = window.innerWidth - x - X_OFFSET;
      menuElement.style.left = "auto";
      menuElement.style.right = x + "px";
    } else {
      x += X_OFFSET;
      menuElement.style.left = x + "px";
      menuElement.style.right = "auto";
    }
    if (window.innerHeight - e.pageY < h) {
      y -= h;
    }
    menuElement.style.display = 'block';
    menuElement.style.top = y + "px";
  }

  function getCurrentMenuHeight() {
    // Easier than exact calculation, works just as well.
    let height = 18;
    menuItems.forEach(i => i.sep ? height += 17 : height += 30);
    return height;
  }

  function getCurrentMenuWidth() {
    // Easier than exact calculation, works nearly as well.
    return menuItems.reduce((a, b) => a.text?.length > b.text?.length ? a : b).text.length * 7.6;
  }

  function hideMenu() {
    document.getElementById("context-menu").style.display = 'none';
    if (currentClickedElement) {
      currentClickedElement.classList.remove('hover');
    }
    currentClickedElement = null;
    currentMenu = null;
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.addEventListener("contextmenu", function (e: MouseEvent & { target: HTMLElement }) {
      if (currentClickedElement) {
        e.preventDefault();
        hideMenu();
        return;
      }
      // todo: e.preventDefault() for everywhere in production mode

      if (typeof menus[e.target.dataset.menu] !== 'undefined') {
        currentClickedElement = e.target;
        currentMenu = e.target.dataset.menu;
      }
      else {
        let n = e.target.parentNode;
        for (let j = 0; j < 4; j++, n = n.parentNode) {
          if (n instanceof HTMLElement && typeof menus[n.dataset.menu] !== 'undefined') {
            currentClickedElement = n;
            currentMenu = n.dataset.menu;
            break;
          }
        }
      }

      if (currentMenu) {
        e.preventDefault();
        menuItems = menus[currentMenu](currentClickedElement);
        currentClickedElement.classList.add('hover');
        displayMenu(e);
      }
      else {
        hideMenu();
      }
    });

    document.body.addEventListener("keydown", function (e: KeyboardEvent) {
      if (e.key === "Escape") {
        hideMenu();
      }
    });

    document.addEventListener("click", function(e: MouseEvent) {
      hideMenu();
    });
  });
</script>

<div class="context-menu" id="context-menu">
  <ul class="context-menu__items">
    {#each menuItems as item}
      {#if item.text}
        <li class="context-menu__item">
          <div class="context-menu__action"
            on:click={(e) => item.callback(currentClickedElement)}
            on:keyup={(e) => item.callback(currentClickedElement)}>
            {item.text}
          </div>
        </li>
      {:else if item.sep}
        <li class="context-menu__sep"></li>
      {/if}
    {/each}
  </ul>
</div>
