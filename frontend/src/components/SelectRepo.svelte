<script lang="ts">
  import { repos, currentRepo } from 'stores/repos';
  import { branchSelect, repoSelect } from 'stores/ui';

  function selectRepo(id: string) {
    currentRepo.set(id);
    repoSelect.set(false);
  }

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    repoSelect.set(!$repoSelect);
    if ($repoSelect) branchSelect.set(false);
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && $repoSelect) {
      repoSelect.set(false);
    }
  });

  repoSelect.subscribe(v => {
    let list = document.getElementById("all-repos");
    let btn = document.getElementById("current-repo");
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

<button class="btn btn-drop sidebar-big-button" id="current-repo" on:click={toggleList}>
  <div class="sidebar-big-button__label">Current Repo:</div>
  <div class="sidebar-big-button__value">{$repos[$currentRepo]?.Name ?? 'none selected'}</div>
</button>

<div id="all-repos" class="sidebar-dropdown">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="overlay" on:click={() => repoSelect.set(false)}></div>
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
