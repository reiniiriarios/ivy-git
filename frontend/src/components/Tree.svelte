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
    Branches: Ref[];
    Tags: Ref[];
    Remotes: Ref[];
    Heads: Ref[];
  }

  interface Ref {
    Hash: string;
    Name: string;
    ShortName: string;
  }

  let commits: Commit[] = [];
  let HEAD: Ref;
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
            <th class="b">Branch/Tag</th>
            <th>Tree</th>
            <th>Commit</th>
            <th>Author</th>
            <th>Date</th>
          </tr>
          {#each Object.entries(commits) as [_, commit]}
            <tr>
              <td>
                <div class="tree__refs c{currentColor}">

                  {#if commit.Branches && commit.Branches.length}
                    {#each commit.Branches as b}
                      <div class="tree__branch">
                        <div class="tree__branch-name">{b.Name}</div>
                        {#if commit.Remotes && commit.Remotes.length}
                          {#each commit.Remotes as r}
                            <div class="tree__leaf">{r.ShortName}</div>
                          {/each}
                        {/if}
                      </div>
                    {/each}
                  {:else if commit.Remotes && commit.Remotes.length}
                    {#each commit.Remotes as r}
                      <div class="tree__branch">
                        <div class="tree__leaf">{r.Name}</div>
                      </div>
                    {/each}
                  {/if}

                  {#if commit.Hash == HEAD.Hash}
                    <div class="tree__branch">
                      <div class="tree__branch-name">HEAD</div>
                      {#if commit.Heads && commit.Heads.length}
                        {#each commit.Heads as h}
                          <div class="tree__leaf">{h.ShortName}</div>
                        {/each}
                      {/if}
                    </div>
                  {:else if commit.Heads && commit.Heads.length}
                    {#each commit.Heads as h}
                      <div class="tree__branch">
                        <div class="tree__leaf">{h.Name}</div>
                      </div>
                    {/each}
                  {/if}

                  {#if commit.Tags && commit.Tags.length}
                    {#each commit.Tags as t}
                      <div class="tree__tag">{t.Name}</div>
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
        text-align: left;
        padding: 0.25rem 0.5rem;
        height: 2rem;
        white-space: nowrap;
        background-color: var(--color-scale-gray-7);

        &:not(:first-child) {
          border-left: 1px solid var(--color-scale-gray-8);
        }

        &.b {
          text-align: right;
          padding-right: 1rem;
        }
      }

      td {
        text-align: left;
        padding: 0.125rem 0.5rem;
        height: 1.75rem;
        white-space: nowrap;
      }
    }

    &__refs {
      display: flex;
      justify-content: right;
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
        width: 0.5rem;
        position: absolute;
        right: -0.5rem;
        background-color: red;
      }

      &:last-child::after {
        width: 1rem;
        right: -1rem;
      }

      &-name {
        padding: 0.2rem 0.5rem;
      }
    }

    &__leaf {
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
