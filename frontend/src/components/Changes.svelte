<script lang="ts">
  import { GitListChanges } from "../../wailsjs/go/main/App";

  interface Change {
    File: string;
    Basename: string;
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
          <span class="changes__status changes__status--{change.Flag}" aria-label="{change.Letter}"></span>
          <span class="changes__file">{change.Basename}</span>
        </li>
      {/each}
    </ul>
  {/if}
  {#if changesY.length}
    <div class="changes__header">Changes</div>
    <ul class="changes__list changes__list--y">
      {#each changesY as change}
        <li class="changes__change">
          <span class="changes__status changes__status--{change.Flag}" aria-label="{change.Letter}"></span>
          <span class="changes__file">{change.Basename}</span>
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
      text-transform: uppercase;
      font-size: 0.8rem;
      color: var(--color-text-label);
      padding: 0.5rem 0.75rem 0.25rem;
      border-bottom: 1px solid var(--color-changes-border);
    }

    &__list {
      list-style: none;
      margin: 0 0 0.5rem;
      padding: 0;
    background-color: var(--color-changes-list-bg);

      li {
        border-bottom: 1px solid var(--color-changes-border);
        padding: 0.5rem 0.75rem 0.25rem;
      }

      &--x {

      }

      &--y {

      }
    }

    &__change {
      display: flex;
    }

    &__status {
      width: 1rem;
      height: 1rem;
      margin-right: 0.5rem;
      position: relative;

      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 1rem;
        height: 1rem;
        text-align: center;
      }

      &--added {
        &::after {
          content: 'A';
          color: var(--color-green-700);
        }
      }

      &--copied {
        &::after {
          content: 'C';
          color: var(--color-lime-800);
        }
      }

      &--deleted {
        &::after {
          content: 'D';
          color: var(--color-red-700);
        }
      }

      &--modified {
        &::after {
          content: 'M';
          color: var(--color-amber-800);
        }
      }

      &--renamed {
        &::after {
          content: 'R';
          color: var(--color-cyan-500);
        }
      }

      &--type-changed {
        &::after {
          content: 'T';
          color: var(--color-teal-900);
        }
      }

      &--unmerged {
        &::after {
          content: 'U';
          color: var(--color-purple-900);
        }
      }

      &--untracked {
        &::after {
          content: 'U';
          color: var(--color-green-300);
        }
      }

      &--ignored {
        &::after {
          content: '!';
          color: var(--color-bluegrey-900);
        }
      }

      &--not-updated {
        &::after {
          content: '-';
          color: var(--color-grey-900);
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

    &__file {

    }
  }
</style>
