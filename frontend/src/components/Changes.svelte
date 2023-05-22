<script lang="ts">
  import { changes } from 'stores/changes';
  import { newCommit } from 'stores/new-commit';

  // Tick states: '', 'checked', 'partial'.
  function tick(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    let set = e.currentTarget.dataset.set;
    let state = e.currentTarget.dataset.state;
    let file = e.currentTarget.dataset.file;
    let newState = state === 'checked' ? '' : 'checked';
    newCommit.setFile(set, file, newState);
  }
</script>

<div class="changes">
  {#if $changes.x.length}
    <div class="changes__header">Staged</div>
    <ul class="changes__list changes__list--x">
      {#each $changes.x as change}
        <li class="change" data-file="{change.File}" data-set="X" data-state="{$newCommit.X.Files[change.File] ?? 'checked'}" on:click={tick} on:keypress={tick}>
          <span class="change__tick">
            <span class="change__tick-box"></span>
            <input type="hidden" name="change[{change.File}]" value="0">
          </span>
          <span class="change__file">
            <span class="change__filename">{change.Basename}</span>
            <span class="change__dir">{change.Dir != '.' ? change.Dir : ''}</span>
          </span>
          <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
        </li>
      {/each}
    </ul>
  {/if}
  {#if $changes.y.length}
    <div class="changes__header">Changes</div>
    <ul class="changes__list changes__list--y">
      {#each $changes.y as change}
        <li class="change" data-file="{change.File}" data-set="Y" data-state="{$newCommit.Y.Files[change.File] ?? 'checked'}" on:click={tick} on:keypress={tick}>
          <span class="change__tick">
            <span class="change__tick-box"></span>
          </span>
          <div class="change__file">
            <div class="change__filename">{change.Basename}</div>
            <div class="change__dir">{change.Dir != '.' ? change.Dir : ''}</div>
          </div>
          <span class="change__status change__status--{change.Flag}" aria-label="{change.Letter}"></span>
        </li>
      {/each}
    </ul>
  {/if}
</div>
