<script lang="ts">
  import { CreateBranch } from 'src/_tmp';

  import { checkRef } from 'scripts/check-ref';
  import { parseResponse } from 'scripts/parse-response';

  import { branches, currentBranch, type Branch } from 'stores/branches';
  import { commitData, commitSignData } from 'stores/commit-data';
  import { messageDialog } from 'stores/message-dialog';
  import { branchSelect, repoSelect } from 'stores/ui';

  function newBranch() {
    messageDialog.confirm({
      heading: 'Create Branch',
      message: `Create a branch?`,
      blank: "Name of Branch",
      validateBlank: checkRef,
      confirm: 'Create',
      callbackConfirm: () => {
        CreateBranch(messageDialog.blankValue(), "", true).then(r => {
          parseResponse(r, () => {
            currentBranch.set({Name: messageDialog.blankValue()} as Branch);
            branchSelect.set(false);
            commitData.refresh();
            commitSignData.refresh();
          })
        });
      }
    });
  }

  function switchBranch(b: string) {
    currentBranch.switch(b);
    branchSelect.set(false);
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
  class:detached={$currentBranch.Name === 'HEAD'}
  style:display={$repoSelect ? 'none' : 'flex'}
  on:click={toggleList}
>
  <div class="sidebar-big-button__label">Current Branch:</div>
  <div class="sidebar-big-button__value">{
    $currentBranch?.Name
      ? $currentBranch.Name === 'HEAD'
        ? 'DETACHED HEAD'
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
