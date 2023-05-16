<script lang="ts">
  import { branches, currentBranch } from 'stores/branches';
  import { branchSelect, repoSelect } from 'stores/ui';

  function newBranch() {
    // ...
  }

  function switchBranch(b: string) {
    currentBranch.set(b);
    branchSelect.set(false);
  }

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    branchSelect.set(!$branchSelect);
    if ($branchSelect) repoSelect.set(false);
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && $branchSelect) {
      branchSelect.set(false);
    }
  });

  branchSelect.subscribe(v => {
    let list = document.getElementById("all-branches");
    let btn = document.getElementById("current-branch");
    if (list && btn) {
      if (v) {
        list.style.display = "block";
        btn.classList.add("active");
      } else {
        list.style.display = "none";
        btn.classList.remove("active");
      }
    }
  });
</script>

<button class="btn btn-drop sidebar-big-button" id="current-branch" on:click={toggleList} class:active={$branchSelect}>
  <div class="sidebar-big-button__label">Current Branch:</div>
  <div class="sidebar-big-button__value">{$currentBranch?.Name ?? "none selected"}</div>
</button>
<div id="all-branches" class="sidebar-dropdown" style:display={$branchSelect ? 'block' : 'none'}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="overlay" on:click={() => branchSelect.set(false)} />
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
