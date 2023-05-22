<script lang="ts">
  import octicons from '@primer/octicons';
  import { parseResponse } from 'scripts/parse-response';
  import { changes } from 'stores/changes';
  import { currentFile } from 'stores/current-file';
  import { currentTab } from 'stores/current-tab';
  import { StageFile, UnstageFile, StageAll, UnstageAll } from 'wailsjs/go/main/App'

  function stage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    StageFile(e.currentTarget.dataset.file).then(result => {
      parseResponse(result, () => {
        changes.refresh();
      });
    });
  }

  function unstage(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    UnstageFile(e.currentTarget.dataset.file).then(result => {
      parseResponse(result, () => {
        changes.refresh();
      });
    });
  }

  function stageAll() {
    StageAll().then(result => {
      parseResponse(result, () => {
        changes.refresh();
      });
    });
  }

  function unstageAll() {
    UnstageAll().then(result => {
      parseResponse(result, () => {
        changes.refresh();
      });
    });
  }

  function selectFile(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    currentTab.set('changes');
    currentFile.select(e.currentTarget.dataset.file, e.currentTarget.dataset.staged === 'true');
  }
</script>

<div class="changes">
  {#if $changes.x.length}
    <div class="changes__header">
      <div class="changes__header-section">Staged</div>
      <div class="changes__header-stage-all changes__header-stage-all--unstage">
        <button class="textbtn" on:click={unstageAll}>All {@html octicons['arrow-down'].toSVG({width: 14})}</button>
      </div>
    </div>
    <ul class="changes__list changes__list--x">
      {#each $changes.x as change}
        <li class="change" class:active={$currentFile.File === change.File && $currentFile.Staged}>
          <span class="change__stage change__stage--unstage"
            aria-label="Unstage File"
            data-file="{change.File}"
            on:click={unstage}
            on:keypress={unstage}>
            {@html octicons['arrow-down'].toSVG({width: 16})}
          </span>
          <div class="change__file"
            data-file="{change.File}"
            data-staged="true"
            on:click={selectFile}
            on:keypress={selectFile}>
            <span class="change__filename">
              <span class="change__basename">{change.Basename}</span>
              <span class="change__dir">{change.Dir != '.' ? change.Dir : ''}</span>
            </span>
            <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
  {#if $changes.y.length}
    <div class="changes__header">
      <div class="changes__header-section">Changes</div>
      <div class="changes__header-stage-all changes__header-stage-all--stage">
        <button class="textbtn" on:click={stageAll}>All {@html octicons['arrow-up'].toSVG({width: 14})}</button>
      </div>
    </div>
    <ul class="changes__list changes__list--y">
      {#each $changes.y as change}
        <li class="change" class:active={$currentFile.File === change.File && !$currentFile.Staged}>
          <span class="change__stage"
            aria-label="Stage File"
            data-file="{change.File}"
            on:click={stage}
            on:keypress={stage}>
            {@html octicons['arrow-up'].toSVG({width: 16})}
          </span>
          <div class="change__file"
            data-file="{change.File}"
            data-staged="false"
            on:click={selectFile}
            on:keypress={selectFile}>
            <div class="change__filename">
              <div class="change__basename">{change.Basename}</div>
              <div class="change__dir">{change.Dir != '.' ? change.Dir : ''}</div>
            </div>
            <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>
