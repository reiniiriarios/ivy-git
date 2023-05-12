<script lang="ts">
  import { changes } from 'stores/changes';
</script>

<div class="changes">
  {#if $changes.x.length}
    <div class="changes__header">Staged</div>
    <ul class="changes__list changes__list--x">
      {#each $changes.x as change}
        <li class="changes__change">
          <span class="changes__file">
            <span class="changes__filename">{change.Basename}</span>
            <span class="changes__dir">{change.Dir != '.' ? change.Dir : ''}</span>
          </span>
          <span class="changes__status changes__status--{change.Flag}" aria-label="{change.Letter}"></span>
        </li>
      {/each}
    </ul>
  {/if}
  {#if $changes.y.length}
    <div class="changes__header">Changes</div>
    <ul class="changes__list changes__list--y">
      {#each $changes.y as change}
        <li class="changes__change">
          <div class="changes__file">
            <div class="changes__filename">{change.Basename}</div>
            <div class="changes__dir">{change.Dir != '.' ? change.Dir : ''}</div>
          </div>
          <span class="changes__status changes__status--{change.Flag}" aria-label="{change.Letter}"></span>
        </li>
      {/each}
    </ul>
  {/if}
</div>
