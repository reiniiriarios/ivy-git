<script lang="ts">
  import { GitListChanges } from 'wailsjs/go/main/App';

  interface Change {
    File: string;
    Basename: string;
    Dir: string;
    Letter: string;
    Flag: string;
  }

  let changesX: Change[] = [];
  let changesY: Change[] = [];

  (window as any).getChanges = () => {
    GitListChanges().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          changesX = result.ChangesX ?? [];
          changesY = result.ChangesY ?? [];
          break;
      }
    });
  };
</script>

<div class="changes">
  {#if changesX.length}
    <div class="changes__header">Staged</div>
    <ul class="changes__list changes__list--x">
      {#each changesX as change}
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
  {#if changesY.length}
    <div class="changes__header">Changes</div>
    <ul class="changes__list changes__list--y">
      {#each changesY as change}
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
