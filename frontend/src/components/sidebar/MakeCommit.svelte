<script lang="ts">
  import octicons from "@primer/octicons";
  import { isDarwin } from "scripts/env";
  import { parseResponse } from "scripts/parse-response";
  import { currentBranch, detachedHead } from "stores/branches";
  import { changes, changesNumConflicts, changesNumStaged, changesNumUnstaged, mergeConflicts, mergeConflictsResolved } from "stores/changes";
  import { commitData } from "stores/commits";
  import { repoState } from "stores/repo-state";
  import { repoSelect, branchSelect, commitMessageSubject, commitMessageBody, commitMessage } from "stores/ui";
  import { MakeCommit, MakeStash } from "wailsjs/go/main/App";

  let running: boolean = false;
  let dropdownOpen: boolean = false;
  let state: string = 'commit';

  function stateCommit() {
    if (state === 'amend') {
      $commitMessageSubject = '';
      $commitMessageBody = '';
    }
    state = 'commit';
  }

  function stateAmend() {
    commitMessage.fetchLast();
    state = 'amend';
  }

  function stateStash() {
    if (state === 'amend') {
      $commitMessageSubject = '';
      $commitMessageBody = '';
    }
    state = 'stash';
  }

  function make() {
    if ($commitMessageSubject && !running && ($changesNumStaged || $changesNumUnstaged)) {
      running = true;
      if (state === 'commit' || state === 'amend') {
        MakeCommit($commitMessageSubject, $commitMessageBody, state === 'amend').then(result => {
          parseResponse(result, () => {
            commitMessage.clear();
            changes.refresh();
            commitData.refresh();
            state = 'commit';
          });
          running = false;
        });
      }
      else if (state === 'stash') {
        MakeStash($commitMessageSubject).then(result => {
          parseResponse(result, () => {
            commitMessage.clear();
            changes.refresh();
            commitData.refresh();
            state = 'commit';
          });
          running = false;
        });
      }
      else {
        console.error('Unrecognized commit button state');
      }
    }
  }

  function makeCommitKeypress(e: KeyboardEvent & { currentTarget: HTMLElement }) {
    let cmd = (isDarwin() && e.metaKey) || (!isDarwin() && e.ctrlKey);
    if (cmd && (e.key === '\n' || e.key === 'Enter')) {
      make();
      e.currentTarget.blur();
    }
  }

  document.addEventListener("click", function(e: MouseEvent & { target: HTMLElement }) {
    // HACK: closest() isn't correctly detecting anything outside of the svg if the svg is clicked.
    // This should be fixed later.
    if (!e.target.closest('.btn__arrow') && !e.target.closest('.octicon')) {
      dropdownOpen = false;
    }
  });
</script>

<div
  class="make-commit"
  class:detached={$detachedHead}
  style:display={$repoSelect || $branchSelect ? 'none' : 'block'}
>
  <div class="make-commit__subject">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Summary</span>
      <input class="text-input repo-state--{$repoState}" type="text" bind:value={$commitMessageSubject} on:keypress={makeCommitKeypress}>
    </label>
  </div>
  <div class="make-commit__body">
    <label class="make-commit__label">
      <span class="make-commit__label-desc">Description</span>
      <textarea class="text-input repo-state--{$repoState}" disabled={state === 'stash'} bind:value={$commitMessageBody} on:keypress={makeCommitKeypress}></textarea>
    </label>
  </div>
  <div class="make-commit__button repo-state--{$repoState} make-commit__button--{state}">
    <div class="btn-arrow" class:btn-arrow--open={dropdownOpen}>
      <button
        class="btn btn__main"
        disabled={
          !$commitMessageSubject
          || running
          || (!$changesNumStaged && !$changesNumUnstaged && !$changesNumConflicts)
          || ($mergeConflicts && !$mergeConflictsResolved)
        }
        on:click={make}
      >
        {#if state === 'commit'}
          Commit to
        {:else if state === 'amend'}
          Amend commit on
        {:else if state === 'stash'}
          Stash on
        {:else}
          Commit to
        {/if}
        <strong>{$currentBranch?.Name}</strong>
      </button>
      <button
        class="btn btn__arrow"
        aria-label="Commit Options"
        on:click={() => dropdownOpen = !dropdownOpen}
        class:btn__arrow--faux-disabled={
          !$commitMessageSubject
          || running
          || (!$changesNumStaged && !$changesNumUnstaged && !$changesNumConflicts)
          || ($mergeConflicts && !$mergeConflictsResolved)
        }
      >
        {#if dropdownOpen}
          {@html octicons["triangle-up"].toSVG({width: 18})}
        {:else}
          {@html octicons["triangle-down"].toSVG({width: 18})}
        {/if}
      </button>
      <div class="btn-arrow__optionscontainer">
        <ul class="btn-arrow__options">
          <li class="btn-arrow__option" on:click={stateCommit} on:keypress={stateCommit}>
            New Commit
          </li>
          <li class="btn-arrow__option" on:click={stateAmend} on:keypress={stateAmend}>
            Amend Commit
          </li>
          <li class="btn-arrow__option" on:click={stateStash} on:keypress={stateStash}>
            Stash
          </li>
        </ul>
      </div>
    </div>
  </div>
</div>
