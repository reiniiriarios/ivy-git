<script lang="ts">
  import { UNCOMMITED_HASH, type Commit } from "../scripts/graph";

  export let commit: Commit;
</script>

<tr class="c-{commit.Color} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}">
  <td><span class="bg"></span>{commit.Subject}</td>
  <td>{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td>{commit.AuthorDatetime}</td>
</tr>

<style lang="scss">
  tr {
    cursor: default !important;
    user-select: none;
    -webkit-user-select: none;
    position: relative;

    &:hover {
      background-color: var(--color-commitlistitem-bg-hover);

      .bg {
        box-sizing: border-box;
        display: block;
        position: absolute;
        top: 0;
        height: var(--commit-details-height);
        background: var(--color-commitlistitem-bg-hover-fade);
        left: -6rem;
        width: 6rem;
        pointer-events: none;
        z-index: -1;
      }
    }
  }

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
