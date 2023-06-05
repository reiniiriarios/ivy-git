<script lang="ts">
  import octicons from '@primer/octicons';
  import { parseResponse } from 'scripts/parse-response';
  import { changes, mergeConflicts } from 'stores/changes';
  import { currentDiff } from 'stores/diffs';
  import { currentTab, branchSelect, repoSelect } from 'stores/ui';
  import { StageFile, UnstageFile, StageAll, UnstageAll, StagePartialFile, UnstagePartialFile } from 'wailsjs/go/main/App'

  function stage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    if (e.currentTarget.dataset.partial === 'true') {
      let f = $changes.y[e.currentTarget.dataset.file];
      StagePartialFile(f.Diff, f.File, f.Letter).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          if ($currentDiff.File === f.File && !$currentDiff.Staged && $currentTab === 'changes') {
            currentDiff.refresh();
          }
        });
      });
    } else {
      StageFile(e.currentTarget.dataset.file).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          if ($currentDiff.File === e.currentTarget.dataset.file && !$currentDiff.Staged && $currentTab === 'changes') {
            currentDiff.clear();
          }
        });
      });
    }
  }

  function unstage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    if (e.currentTarget.dataset.partial === 'true') {
      let f = $changes.x[e.currentTarget.dataset.file];
      UnstagePartialFile(f.Diff, f.File, f.Letter).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          if ($currentDiff.File === f.File && $currentDiff.Staged && $currentTab === 'changes') {
            currentDiff.refresh();
          }
        });
      });
    } else {
      UnstageFile(e.currentTarget.dataset.file).then(result => {
        parseResponse(result, () => {
          changes.refresh();
          if ($currentDiff.File === e.currentTarget.dataset.file && $currentDiff.Staged && $currentTab === 'changes') {
            currentDiff.clear();
          }
        });
      });
    }
  }

  function stageAll() {
    StageAll().then(result => {
      parseResponse(result, () => {
        changes.refresh();
        if ($currentTab === 'changes') {
          if ($currentDiff.Staged) {
            currentDiff.refresh();
          } else {
            currentDiff.clear();
          }
        }
      });
    });
  }

  function unstageAll() {
    UnstageAll().then(result => {
      parseResponse(result, () => {
        changes.refresh();
        if ($currentTab === 'changes') {
          if ($currentDiff.Staged) {
            currentDiff.clear();
          } else {
            currentDiff.refresh();
          }
        }
      });
    });
  }

  function selectFile(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    currentTab.set('changes');
    changes.fetchDiff(e.currentTarget.dataset.list, e.currentTarget.dataset.file);
  }
</script>

<div class="changes" style:display={$repoSelect || $branchSelect ? 'none' : 'block'}>
  {#if $mergeConflicts}
    <div class="changes__header">
      <div class="changes__header-section">Conflicts</div>
    </div>
    <ul class="changes__list changes__list--conflicts">
      {#each Object.entries($changes.c) as [_, change]}
        <li
          class="change"
          class:change--active={$currentDiff.File === change.File && $currentDiff.Conflict && $currentTab === 'changes'}
          class:change--partial={change.Diff?.SelectableLines !== change.Diff?.SelectedLines}
          class:change--none={change.Diff?.SelectedLines === 0}
          class:change--unresolved={!change.Diff?.Resolved}
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
          class:change--none={change.Diff?.SelectedLines === 0}
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
          class:change--none={change.Diff?.SelectedLines === 0}
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
</div>
