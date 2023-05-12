<script lang="ts">
  import { branches, currentBranch } from "../../src/stores/branches";

  let listVisible: boolean = false;

  function newBranch() {
    // ...
  }

  function switchBranch(b: string) {
    currentBranch.set(b);
    hideList();
  }

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    if (listVisible) {
      hideList();
    } else {
      showList();
    }
  }

  function showList() {
    document.getElementById("all-branches").style.display = "block";
    document.getElementById("current-branch").classList.add("active");
    listVisible = true;
  }

  function hideList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    document.getElementById("all-branches").style.display = "none";
    document.getElementById("current-branch").classList.remove("active");
    listVisible = false;
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && listVisible) {
      hideList();
    }
  });
</script>

<button class="btn btn-drop sidebar-big-button" id="current-branch" on:click={toggleList}>
  <div class="sidebar-big-button__label">Current Branch:</div>
  <div class="sidebar-big-button__value">{$currentBranch?.Name ?? "none selected"}</div>
</button>

<div id="all-branches" class="sidebar-dropdown">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="overlay" on:click={hideList} />
  <div class="sidebar-dropdown__container">
    <div class="sidebar-dropdown__bar">
      <div class="sidebar-dropdown__add">
        <button class="btn" on:click={newBranch}>Create Branch +</button>
      </div>
      <ul class="sidebar-dropdown__list">
        {#each Object.entries($branches) as [_, branch]}
          <li>
            <button class="list-btn name" on:click={() => switchBranch(branch?.Name)}>{branch?.Name}</button>
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>
