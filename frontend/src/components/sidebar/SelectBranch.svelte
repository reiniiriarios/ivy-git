<script lang="ts">
  import { branches, currentBranch, detachedHead } from 'stores/branches';
  import { branchSelect, repoSelect } from 'stores/ui';
  import createBranch from 'actions/branch/create';
  import { mergeConflicts } from 'stores/changes';
  import deleteBranch from 'actions/branch/delete';
  import switchBranch from 'actions/branch/switch';
  import { currentRepo, repos } from 'stores/repos';

  const newBranch = () => createBranch();

  function toggleList(e?: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', '\n', 'Enter'].includes(e.key)) {
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
</script>

<button
  id="current-branch"
  class="btn btn-drop sidebar-big-button"
  class:active={$branchSelect}
  class:detached={$detachedHead}
  style:display={$repoSelect ? 'none' : 'flex'}
  data-menu="branchList"
  on:click={toggleList}
>
  <div class="sidebar-big-button__label">Current Branch:</div>
  <div class="sidebar-big-button__value">{
    $currentBranch?.Name
      ? $detachedHead
        ? 'DETACHED HEAD'
        : $mergeConflicts
          ? $currentBranch.Name + ' (Merge Conflicts)'
          : $currentBranch.Name
      : "none selected"
  }</div>
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
        {#if $branches?.length}
          {#each Object.entries($branches) as [_, branch]}
            <li
              class="sidebar-dropdown__item"
              class:sidebar-dropdown__item--selected={branch?.Name === $currentBranch.Name}
              data-menu="branchInList"
              data-name="{branch.Name}"
              data-upstream="{branch.Upstream}"
              data-current="{$currentBranch.Name === branch.Name}"
            >
              <button class="list-btn name" on:click={() => switchBranch(branch?.Name)}>{branch?.Name}</button>
              {#if branch?.Name && branch.Name !== $currentBranch.Name && branch.Name !== $repos[$currentRepo].Main }
                <button class="list-btn x" on:click={() => deleteBranch(branch.Name, !!branch.Upstream)}>&times;</button>
              {/if}
            </li>
          {/each}
        {/if}
      </ul>
      <div class="sidebar-dropdown__remaining-space" data-menu="branchList"></div>
    </div>
  </div>
</div>
