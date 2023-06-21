<script lang="ts">
  import octicons from '@primer/octicons';
  import { messageDialog } from 'stores/message-dialog';
  import { repos, currentRepo } from 'stores/repos';
  import { branchSelect, repoSelect } from 'stores/ui';

  function selectRepo(id: string) {
    currentRepo.switch(id);
    repoSelect.set(false);
  }

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', '\n', 'Enter'].includes(e.key)) {
      return;
    }
    repoSelect.set(!$repoSelect);
    if ($repoSelect) branchSelect.set(false);
  }

  function addRepo() {
    messageDialog.options({
      heading: 'Add Repo',
      options: [
        {
          text: 'Add Existing Repo',
          icon: octicons['project-symlink'].toSVG({ width: 32 }),
          callback: repos.add,
        },
        {
          text: 'Clone Repo',
          icon: octicons['repo-clone'].toSVG({ width: 32 }),
          callback: repos.clone,
        },
        {
          text: 'Create New Repo',
          icon: octicons['plus-circle'].toSVG({ width: 32 }),
          callback: repos.create,
        },
      ],
    });
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && $repoSelect) {
      repoSelect.set(false);
    }
  });
</script>

<button class="btn btn-drop sidebar-big-button" id="current-repo" on:click={toggleList} class:active={$repoSelect}>
  <div class="sidebar-big-button__label">Current Repo:</div>
  <div class="sidebar-big-button__value">{$repos[$currentRepo]?.Name ? $repos[$currentRepo]?.Name : 'none selected'}</div>
</button>

<div id="all-repos" class="sidebar-dropdown" style:display={$repoSelect ? 'block' : 'none'}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="overlay" on:click={() => repoSelect.set(false)}></div>
  <div class="sidebar-dropdown__container">
    <div class="sidebar-dropdown__bar">
      <div class="sidebar-dropdown__add">
        <button class="btn" on:click={addRepo}>Add Repo +</button>
      </div>
      <ul class="sidebar-dropdown__list">
        {#each Object.entries($repos) as [id, repo]}
          {#if repo?.Name}
            <li class="sidebar-dropdown__item" class:sidebar-dropdown__item--selected={id === $currentRepo}>
              <button class="list-btn name" on:click={() => selectRepo(id)}>{repo.Name}</button>
              <button class="list-btn x" on:click={() => repos.delete(id)}>&times;</button>
            </li>
          {/if}
        {/each}
      </ul>
    </div>
  </div>
</div>
