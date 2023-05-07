<script lang="ts">
  import { UNCOMMITED_HASH, type Commit } from "../scripts/graph";

  export let commit: Commit;
</script>

<tr class="c-{commit.Color} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}">
  <td>{commit.Subject}</td>
  <td>{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td>{commit.AuthorDatetime}</td>
</tr>

<style lang="scss">
  td {
    text-align: left;
    white-space: nowrap;
    padding-left: 0.67rem;
    padding-right: 0.67rem;
    height: var(--commit-details-height);
    box-sizing: border-box;

    &:last-child {
      padding-right: 0.75rem;
    }
  }

  .uncommitted {
    color: var(--color-scale-a-3-100);
  }

  .merge {
    color: var(--color-scale-a-3-100);
  }

  .stash {
    color: var(--color-scale-a-2-100);
  }
</style>
