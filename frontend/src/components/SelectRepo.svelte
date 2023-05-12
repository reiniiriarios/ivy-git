<script lang="ts">
  import { repos, currentRepo } from '../stores/repos';

  let listVisible: boolean = false;

  function selectRepo(id: string) {
    currentRepo.set(id);
    hideList();
  }

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    listVisible ? hideList() : showList();
  }

  function showList() {
    document.getElementById("all-repos").style.display = "block";
    document.getElementById("current-repo").classList.add('active');
    listVisible = true;
  }

  function hideList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    document.getElementById("all-repos").style.display = "none";
    document.getElementById("current-repo").classList.remove('active');
    listVisible = false;
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && listVisible) {
      hideList();
    }
  });
</script>

<button class="btn btn-drop sidebar-big-button" id="current-repo" on:click={toggleList}>
  <div class="sidebar-big-button__label">Current Repo:</div>
  <div class="sidebar-big-button__value">{$repos[$currentRepo]?.Name ?? 'none selected'}</div>
</button>

<div id="all-repos" class="sidebar-dropdown">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="overlay" on:click={hideList}></div>
  <div class="sidebar-dropdown__container">
    <div class="sidebar-dropdown__bar">
      <div class="sidebar-dropdown__add">
        <button class="btn" on:click={repos.add}>Add Repo +</button>
      </div>
      <ul class="sidebar-dropdown__list">
        {#each Object.entries($repos) as [id, repo]}
          <li>
            <button class="list-btn name" on:click={() => selectRepo(id)}>{repo.Name}</button>
            <button class="list-btn x" on:click={() => repos.delete(id)}>&times;</button>
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>
