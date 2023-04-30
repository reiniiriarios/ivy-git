<script lang="ts">
  export let active: boolean;

  import { GetCommitsForTree } from "../../wailsjs/go/main/App";

  interface Commit {
    Hash: string;
    Parents: string[];
    AuthorName: string;
    AuthorEmail: string;
    AuthorTimestamp: number;
    AuthorDatetime: string;
    Subject: string;
    Tags: string[];
    Remotes: string[];
    Heads: string[];
  }

  let commits: Commit[] = [];
  let HEAD: string = "";
  let currentColor = 1;

  (window as any).GetCommitsForTree = () => {
    GetCommitsForTree().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          commits = result.Commits as Commit[];
          HEAD = result.HEAD;
          console.log(commits);
          break;
      }
    });
  };
</script>

{#if active}
  <div class="tree">
    <div class="tree__table">
      {#if Object.entries(commits).length}
        <table>
          <tr>
            <th>Branch/Tag</th>
            <th>Tree</th>
            <th>Commit</th>
            <th>Author</th>
            <th>Date</th>
          </tr>
          {#each Object.entries(commits) as [_, commit]}
            <tr>
              <td>
                <div class="tree__refs c{currentColor}">
                  {#if commit.Heads && commit.Heads.length}
                    {#each commit.Heads as h}
                      <div class="tree__branch">
                        <div class="tree__branch-name">{h}</div>
                        {#if commit.Remotes && commit.Remotes.length}
                          {#each commit.Remotes as r}
                            <div class="tree__remote">{r}</div>
                          {/each}
                        {/if}
                        {#if commit.Hash == HEAD}
                          <div class="tree__head">HEAD</div>
                        {/if}
                      </div>
                    {/each}
                  {/if}
                  {#if commit.Tags && commit.Tags.length}
                    {#each commit.Tags as t}
                      <div class="tree__tag">{t}</div>
                    {/each}
                  {/if}
                </div>
              </td>
              <td>...</td>
              <td>{commit.Subject}</td>
              <td>{commit.AuthorName ?? commit.AuthorEmail}</td>
              <td>{commit.AuthorDatetime}</td>
            </tr>
          {/each}
        </table>
      {/if}
    </div>
  </div>
{/if}

<style lang="scss">
  .tree {
    width: 100%;
    height: 100%;

    &__table {
      width: 100%;
      overflow-x: auto;

      table {
        width: 100%;
      }

      th {
        background-color: var(--color-scale-gray-7);

        &:not(:first-child) {
          border-left: 1px solid var(--color-scale-gray-8);
        }
      }

      th, td {
        text-align: left;
        padding: 0.125rem 0.5rem;
        height: 1.75rem;
        white-space: nowrap;
      }
    }

    &__refs {
      text-align: right;
    }

    &__branch {
      display: inline-flex;
      justify-content: left;
      align-items: center;
      border-left: 2px solid red;
      background-color: rgba(255 0 0 / 25%);
      position: relative;
      margin-right: 0.5rem;

      &::after {
        content: '';
        height: 1px;
        width: 1rem;
        position: absolute;
        right: -1rem;
        background-color: red;
      }

      &-name {
        padding: 0.2rem 0.5rem;
      }
    }

    &__remote {
      padding: 0.2rem 0.5rem;
    }

    &__head {
      padding: 0.2rem 0.5rem;
      background-color: rgba(0 0 0 / 25%);
    }
  }

  .c1 {
    .tree__branch {
      border-color: var(--color-branch-1-border);
      background-color: var(--color-branch-1-bg);

      &::after {
        background-color: var(--color-branch-1-bg);
      }
    }
  }

  .c2 {
    .tree__branch {
      border-color: var(--color-branch-2-border);
      background-color: var(--color-branch-2-bg);

      &::after {
        background-color: var(--color-branch-2-bg);
      }
    }
  }

  .c3 {
    .tree__branch {
      border-color: var(--color-branch-3-border);
      background-color: var(--color-branch-3-bg);

      &::after {
        background-color: var(--color-branch-3-bg);
      }
    }
  }

  .c4 {
    .tree__branch {
      border-color: var(--color-branch-4-border);
      background-color: var(--color-branch-4-bg);

      &::after {
        background-color: var(--color-branch-4-bg);
      }
    }
  }
</style>
