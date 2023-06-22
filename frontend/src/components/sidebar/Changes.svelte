<script lang="ts">
  import octicons, { x } from '@primer/octicons';
  import { stageAll } from 'actions/stage/stage-all';
  import { stageFile } from 'actions/stage/stage-file';
  import { unstageAll } from 'actions/stage/unstage-all';
  import { unstageFile } from 'actions/stage/unstage-file';
  import { changes, mergeConflicts, type Change } from 'stores/changes';
  import { currentDiff } from 'stores/diffs';
  import { currentTab, branchSelect, repoSelect } from 'stores/ui';

  function stage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    let file: Change = $changes.y[e.currentTarget?.dataset?.file];
    if (file) {
      let partial = e.currentTarget?.dataset?.partial === 'true';
      stageFile(file, partial);
    }
  }

  function unstage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    let file: Change = $changes.x[e.currentTarget?.dataset?.file];
    if (file) {
      let partial = e.currentTarget?.dataset?.partial === 'true';
      unstageFile(file, partial);
    }
  }

  function selectFile(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    currentTab.set('changes');
    changes.fetchDiff(e.currentTarget.dataset.list, e.currentTarget.dataset.file);
  }
</script>

<div class="changes" style:display={$repoSelect || $branchSelect ? 'none' : 'flex'}>
  {#if $mergeConflicts}
    <div class="changes__header">
      <div class="changes__header-section">Conflicts</div>
    </div>
    <ul class="changes__list changes__list--conflicts">
      {#each Object.entries($changes.c) as [_, change]}
        <li
          class="change"
          class:change--active={$currentDiff.File === change.File && $currentDiff.Conflict && $currentTab === 'changes'}
          class:change--unresolved={!change.Diff?.Resolved}
          data-menu="change"
          data-file="{change.File}"
          data-conflict="true"
        >
          <div class="change__file"
            data-file="{change.File}"
            data-status="{change.Them}{change.Us}"
            data-conflict="true"
            data-list="c"
            on:click={selectFile}
            on:keypress={selectFile}>
            <span class="change__filename">
              <span class="change__basename">{change.Basename}</span>
              <span class="change__dir">{change.Dir != '.' ? change.Dir : ''}</span>
            </span>
            <span class="change__status change__status--{change.Flag}" aria-label="{change.Them} {change.Us}"></span>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
  {#if Object.keys($changes.x).length}
    <div class="changes__header">
      <div class="changes__header-section">Staged</div>
      <div class="changes__header-stage-all changes__header-stage-all--unstage">
        <button class="textbtn" on:click={unstageAll}>All {@html octicons['arrow-down'].toSVG({width: 14})}</button>
      </div>
    </div>
    <ul class="changes__list changes__list--x">
      {#each Object.entries($changes.x) as [_, change]}
        <li
          class="change"
          class:change--active={$currentDiff.File === change.File && $currentDiff.Staged && $currentTab === 'changes'}
          class:change--partial={change.Diff?.SelectableLines !== change.Diff?.SelectedLines}
          class:change--none={change.Diff?.SelectableLines > 0 && change.Diff?.SelectedLines === 0}
          data-menu="change"
          data-file="{change.File}"
          data-staged="true"
          data-partial="{change.Diff?.SelectableLines !== change.Diff?.SelectedLines}"
        >
          <div class="change__file"
            data-file="{change.File}"
            data-status="{change.Letter}"
            data-staged="true"
            data-list="x"
            on:click={selectFile}
            on:keypress={selectFile}>
            <span class="change__filename">
              <span class="change__basename">{change.Basename}</span>
              <span class="change__dir">{change.Dir != '.' ? change.Dir : ''}</span>
            </span>
            <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
          </div>
          <span class="change__stage change__stage--unstage"
            aria-label="Unstage File"
            data-file="{change.File}"
            data-partial="{change.Diff?.SelectableLines !== change.Diff?.SelectedLines}"
            on:click={unstage}
            on:keypress={unstage}>
            {@html octicons['arrow-down'].toSVG({width: 16})}
          </span>
        </li>
      {/each}
    </ul>
  {/if}
  {#if Object.keys($changes.y).length}
    <div class="changes__header">
      <div class="changes__header-section">Changes</div>
      <div class="changes__header-stage-all changes__header-stage-all--stage">
        <button class="textbtn" on:click={stageAll}>All {@html octicons['arrow-up'].toSVG({width: 14})}</button>
      </div>
    </div>
    <ul class="changes__list changes__list--y">
      {#each Object.entries($changes.y) as [_, change]}
        <li
          class="change"
          class:change--active={$currentDiff.File === change.File && !$currentDiff.Staged && $currentTab === 'changes'}
          class:change--partial={change.Diff?.SelectableLines !== change.Diff?.SelectedLines}
          class:change--none={change.Diff?.SelectableLines > 0 && change.Diff?.SelectedLines === 0}
          data-menu="change"
          data-file="{change.File}"
          data-staged="false"
          data-partial="{change.Diff?.SelectableLines !== change.Diff?.SelectedLines}"
        >
          <div class="change__file"
            data-file="{change.File}"
            data-status="{change.Letter}"
            data-staged="false"
            data-list="y"
            on:click={selectFile}
            on:keypress={selectFile}>
            <div class="change__filename">
              <div class="change__basename">{change.Basename}</div>
              <div class="change__dir">{change.Dir != '.' ? change.Dir : ''}</div>
            </div>
            <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
          </div>
          <span class="change__stage"
            aria-label="Stage File"
            data-file="{change.File}"
            data-partial="{change.Diff?.SelectableLines !== change.Diff?.SelectedLines}"
            on:click={stage}
            on:keypress={stage}>
            {@html octicons['arrow-up'].toSVG({width: 16})}
          </span>
        </li>
      {/each}
    </ul>
  {/if}
  <div class="changes__remaining-space"
    data-menu="changes"
    data-unstaged={$changes?.y ? !!Object.keys($changes.y).length : false}
    data-staged={$changes?.x ? !!Object.keys($changes.x).length : false}
    data-conflicts={$changes?.c ? !!Object.keys($changes.c).length : false}></div>
</div>
