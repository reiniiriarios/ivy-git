<script lang="ts">
  import { GitListChanges } from "../../wailsjs/go/main/App";

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

<style lang="scss">
  .changes {
    text-align: left;
    background-color: var(--color-changes-bg);
    height: 100%;

    &__header {
      font-size: 0.9rem;
      color: var(--color-scale-gray-3);
      padding: 0.5rem 0.75rem 0.5rem;
      border-bottom: 1px solid var(--color-changes-border);
    }

    &__list {
      list-style: none;
      margin: 0 0 0.5rem;
      padding: 0;
      background-color: var(--color-changes-list-bg);
    }

    &__change {
      display: flex;
      justify-content: left;
      align-items: center;
      border-bottom: 1px solid var(--color-changes-border);
      padding: 0.30rem 1.5rem 0.36rem;
      min-height: 3rem;
    }

    &__file {
      display: flex;
      flex-direction: column;
    }

    &__dir {
      color: var(--color-scale-gray-3);
      font-size: 0.85rem;
    }

    &__status {
      position: relative;
      margin-left: auto;

      &::after {
        content: '';
        width: 1.15rem;
        height: 1.15rem;
        font-size: 0.85rem;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      &--added {
        &::after {
          content: 'A';
          color: var(--color-scale-green-3);
        }
      }

      &--copied {
        &::after {
          content: 'C';
          color: var(--color-scale-yellow-7);
        }
      }

      &--deleted {
        &::after {
          content: 'D';
          color: var(--color-scale-red-3);
        }
      }

      &--modified {
        &::after {
          content: 'M';
          color: var(--color-scale-orange-3);
        }
      }

      &--renamed {
        &::after {
          content: 'R';
          color: var(--color-scale-blue-7);
        }
      }

      &--type-changed {
        &::after {
          content: 'T';
          color: var(--color-scale-purple-7);
        }
      }

      &--unmerged {
        &::after {
          content: 'U';
          color: var(--color-scale-pink-7);
        }
      }

      &--untracked {
        &::after {
          content: 'U';
          color: var(--color-scale-green-3);
        }
      }

      &--ignored {
        &::after {
          content: '!';
          color: var(--color-scale-gray-6);
        }
      }

      &--not-updated {
        &::after {
          content: '-';
          color: var(--color-scale-grey-6);
        }
      }

      &--unknown {
        background-color: #00f;

        &::after {
          content: 'X';
          color: #f00;
        }
      }
    }
  }
</style>
