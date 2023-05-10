<script lang="ts">
  import { UNCOMMITED_HASH, type Commit, type Ref } from "../scripts/graph";
  import CommitLabels from "./CommitLabels.svelte";

  export let commit: Commit;
  export let HEAD: Ref;
</script>

<tr data-id="{commit.Id}" class="commit c-{commit.Color} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}">
  <td class="commit__td commit__td--refs">
    {#if commit.Labeled}
      <CommitLabels commit={commit} HEAD={HEAD} />
    {/if}
  </td>
  <td class="commit__td commit__td--tree"></td>
  <td class="commit__td commit__td--subject">{commit.Subject}</td>
  <td class="commit__td commit__td--author">{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td class="commit__td commit__td--authortime">{commit.AuthorDatetime}</td>
</tr>

<!-- svelte-ignore css-unused-selector -->
<style lang="scss">
  .commit {
    cursor: default !important;
    user-select: none;
    -webkit-user-select: none;
    position: relative;

    &__td {
      text-align: left;
      padding-left: 0.67rem;
      padding-right: 0.67rem;
      height: var(--commit-details-height);
      box-sizing: border-box;
      white-space: nowrap;

      &:first-child {
        padding-left: 0.67rem;
        padding-right: 0;
      }

      &:last-child {
        padding-right: 0.75rem;
      }

      &--tree {
        padding: 0;
        overflow: hidden;
      }

      &--subject,
      &--author,
      &--authortime {
        overflow: hidden;
        text-overflow: ellipsis;
      }

      &:not(:first-child) {
        .uncommitted & {
          color: var(--color-scale-a-3-100);
        }

        .merge & {
          color: var(--color-scale-a-3-100);
        }

        .stash & {
          color: var(--color-scale-a-2-100);
        }
      }
    }

    &:hover,
    .hover {
      background-color: var(--color-commitlistitem-bg-hover);
    }
  }
</style>
