<script lang="ts">
  export let active: boolean;

  import octicons from '@primer/octicons';

  import { GetCommitList } from "../../wailsjs/go/main/App";

  const UNCOMMITED_HASH = "#";

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

  (window as any).GetCommitList = () => {
    GetCommitList().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          commits = result.Commits as Commit[];
          HEAD = result.HEAD;
          console.log(HEAD);
          console.log(commits);
          break;
      }
    });
  };
</script>

{#if active}
  <div class="tree">
    {#if Object.entries(commits).length}
      <table>
        <tr>
          <th class="h-b">Branch</th>
          <th>Tree</th>
          <th class="h-c">Commit</th>
          <th>Author</th>
          <th>Date</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <tr class="{commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''}">
            <td>
              <div class="tree__refs c{currentColor}">

                {#if commit.Branches && commit.Branches.length}
                  {#each commit.Branches as b}
                    <div class="tree__branch">
                      <div class="tree__icon">{@html octicons['git-branch'].toSVG()}</div>
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
                      <div class="tree__icon">{@html octicons['git-branch'].toSVG()}</div>
                      <div class="tree__leaf">{r.Name}</div>
                    </div>
                  {/each}
                {/if}

                {#if commit.Hash == HEAD.Hash}
                  <div class="tree__head">
                    <div class="tree__icon">{@html octicons['arrow-right'].toSVG()}</div>
                    <div class="tree__head-name">HEAD</div>
                    {#if commit.Heads && commit.Heads.length}
                      {#each commit.Heads as h}
                        <div class="tree__leaf">{h.ShortName}</div>
                      {/each}
                    {/if}
                  </div>
                {:else if commit.Heads && commit.Heads.length}
                  {#each commit.Heads as h}
                    <div class="tree__head">
                      <div class="tree__icon">{@html octicons['arrow-right'].toSVG()}</div>
                      <div class="tree__leaf">{h.Name}</div>
                    </div>
                  {/each}
                {/if}

                {#if commit.Tags && commit.Tags.length}
                  {#each commit.Tags as t}
                    <div class="tree__tag">
                      <div class="tree__icon">{@html octicons['tag'].toSVG()}</div>
                      <div class="tree__tag-name">{t.Name}</div>
                    </div>
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
{/if}

<style lang="scss">
  .tree {
    min-width: 100%;
    height: calc(100vh - var(--tabs-height) - var(--title-bar-height));
    overflow: auto;

    table {
      min-width: 100%;
      margin-bottom: 0.5rem;

      tr {
        th {
          text-align: left;
          padding: 0.25rem 0.5rem;
          height: 2rem;
          white-space: nowrap;
          background-color: var(--color-scale-gray-7);

          &:not(:first-child) {
            border-left: 1px solid var(--color-scale-gray-8);
          }

          &.h-b {
            text-align: right;
            padding-right: 1rem;
          }

          &.h-c {
            width: 100%;
          }
        }

        td {
          text-align: left;
          padding: 0.125rem 0.5rem;
          height: 1.75rem;
          white-space: nowrap;

          &:first-child {
            padding-left: 0.67rem;
          }

          &:not(:first-child) {
            padding-right: 0.67rem;
          }
        }

        &.uncommitted {
          color: var(--color-scale-gray-3);
        }
      }
    }

    &__refs {
      display: flex;
      justify-content: right;
    }

    &__branch,
    &__head,
    &__tag {
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
        padding: 0.15rem 0.5rem 0.25rem 0.5rem;
        border-left: 1px solid;
        border-color: inherit;
      }

      .c1 & {
        border-color: var(--color-branch-1-border);
        background-color: var(--color-branch-1-bg);

        &::after {
          background-color: var(--color-branch-1-bg);
        }
      }

      .c2 & {
        border-color: var(--color-branch-2-border);
        background-color: var(--color-branch-2-bg);

        &::after {
          background-color: var(--color-branch-2-bg);
        }
      }

      .c3 & {
        border-color: var(--color-branch-3-border);
        background-color: var(--color-branch-3-bg);

        &::after {
          background-color: var(--color-branch-3-bg);
        }
      }

      .c4 & {
        border-color: var(--color-branch-4-border);
        background-color: var(--color-branch-4-bg);

        &::after {
          background-color: var(--color-branch-4-bg);
        }
      }
    }

    &__leaf {
      padding: 0.15rem 0.5rem 0.25rem 0.5rem;
      background-color: rgba(0 0 0 / 25%);
    }

    &__icon {
      width: 1.7rem;
      height: 1.5rem;
      padding-left: 0.1rem;
      padding-right: 0.1rem;
      display: flex;
      justify-content: center;
      align-items: center;
      fill: var(--color-text);
    }
  }
</style>
