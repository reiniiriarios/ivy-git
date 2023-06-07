<script lang="ts">
  import { branches, currentBranch, type Branch, detachedHead } from 'stores/branches';
  import { branchSelect, repoSelect } from 'stores/ui';
  import createBranch from 'actions/create-branch';
  import { mergeConflicts } from 'stores/changes';
  import { RepoState, repoState } from 'stores/repo-state';
  import { messageDialog } from 'stores/message-dialog';

  const newBranch = () => createBranch();

  function switchBranch(b: string) {
    if (![RepoState.Nil, RepoState.None].includes($repoState)) {
      messageDialog.error({
        message: "The repo is currently in a state that you cannot (or should not) switch branches."
      });
    } else {
      currentBranch.switch(b);
      branchSelect.set(false);
    }
  }

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
            <li>
              <button class="list-btn name" on:click={() => switchBranch(branch?.Name)}>{branch?.Name}</button>
            </li>
          {/each}
        {/if}
      </ul>
    </div>
  </div>
</div>
