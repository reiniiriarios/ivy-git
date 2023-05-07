<script lang="ts">
  export let active: boolean;

  import octicons from '@primer/octicons';

  import { GetCommitList } from '../../wailsjs/go/main/App';

  import { drawGraph, getLabelDist, UNCOMMITED_HASH } from '../scripts/graph';
  import type { Commit, Ref } from '../scripts/graph';

  let commits: Commit[] = [];
  let HEAD: Ref;
  let currentColor = 1;
  let svg: SVGSVGElement;

  (window as any).GetCommitList = async () => {
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
          svg = drawGraph(result.Graph);
          console.log(result.Graph);
          break;
      }
    });
  };
</script>

{#if active}
  <div class="commits">
    {#if Object.entries(commits).length}
      <table class="commits__branches">
        <tr>
          <th class="h-b">Branch</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <tr class="commit c-{commit.Color} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''}">
            <td>
              {#if commit.Labeled}
                <div class="commit__refs">

                  {#if commit.Branches && commit.Branches.length}
                    {#each commit.Branches as b}
                      <div class="commit__label commit__branch">
                        <div class="commit__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
                        <div class="commit__label-name">{b.Name}</div>
                        {#if commit.Remotes && commit.Remotes.length}
                          {#each commit.Remotes as r}
                            <div class="commit__leaf">{r.ShortName}</div>
                          {/each}
                        {/if}
                      </div>
                    {/each}
                  {:else if commit.Remotes && commit.Remotes.length}
                    {#each commit.Remotes as r}
                      <div class="commit__label commit__label--branch">
                        <div class="commit__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
                        <div class="commit__leaf">{r.Name}</div>
                      </div>
                    {/each}
                  {/if}

                  {#if commit.Hash == HEAD.Hash}
                    <div class="commit__label commit__label--head">
                      <div class="commit__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
                      <div class="commit__label-name">HEAD</div>
                      {#if commit.Heads && commit.Heads.length}
                        {#each commit.Heads as h}
                          <div class="commit__leaf">{h.ShortName}</div>
                        {/each}
                      {/if}
                    </div>
                  {:else if commit.Heads && commit.Heads.length}
                    {#each commit.Heads as h}
                      <div class="commit__label commit__label--head">
                        <div class="commit__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
                        <div class="commit__leaf">{h.Name}</div>
                      </div>
                    {/each}
                  {/if}

                  {#if commit.Tags && commit.Tags.length}
                    {#each commit.Tags as t}
                      <div class="commit__label commit__label--tag">
                        <div class="commit__icon">{@html octicons['tag'].toSVG({ "width": 14 })}</div>
                        <div class="commit__label-name">{t.Name}</div>
                      </div>
                    {/each}
                  {/if}

                  {#if commit.Stash}
                    <div class="commit__label commit__label--stash">
                      <div class="commit__icon">{@html octicons['inbox'].toSVG({ "width": 14 })}</div>
                      <div class="commit__label-name">{commit.RefName}</div>
                    </div>
                  {/if}

                  <div class="commit__line" style="width:{getLabelDist(commit.X)}px; right:-{getLabelDist(commit.X)}px"></div>
                </div>
              {/if}
            </td>
          </tr>
        {/each}
      </table>
      <div id="tree" class="tree">
        <div class="tree__text">Tree</div>
        <div class="tree__graph">{@html svg.outerHTML}</div>
      </div>
      <table class="commits__details">
        <tr>
          <th class="h-c">Commit</th>
          <th>Author</th>
          <th>Date</th>
        </tr>
        {#each Object.entries(commits) as [_, commit]}
          <tr class="commit c-{currentColor} {commit.Hash === UNCOMMITED_HASH ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}">
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
  .commits {
    min-width: 100%;
    height: calc(100vh - var(--tabs-height) - var(--title-bar-height));
    overflow: auto;
    display: flex;
    flex-direction: row;
    justify-content: stretch;
    align-items: middle;

    table {
      margin-bottom: 0.5rem;

      tr {
        th, td {
          // Height in pixels!
          height: 24px;
          box-sizing: border-box;
        }

        th {
          text-align: left;
          padding: 0.25rem 0.5rem;
          white-space: nowrap;
          background-color: var(--color-scale-a-7-100);

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
          white-space: nowrap;
        }

        &.uncommitted {
          color: var(--color-scale-a-3-100);
        }

        &.merge {
          color: var(--color-scale-a-3-100);
        }

        &.stash {
          color: var(--color-scale-a-2-100);
        }
      }
    }

    &__branches {
      th {
        border-right: 1px solid var(--color-scale-a-8-100);
      }

      td {
        padding-left: 0.67rem;
      }
    }

    &__details {
      flex: 1;

      th {
        border-left: 1px solid var(--color-scale-a-8-100);
      }

      td {
        padding-left: 0.67rem;
        padding-right: 0.67rem;

        &:last-child {
          padding-right: 0.75rem;
        }
      }
    }
  }

  .commit {
    box-sizing: border-box;
    overflow: hidden;

    &__refs {
      display: flex;
      justify-content: right;
      align-items: center;
      position: relative;
    }

    &__label {
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
    }

    &__line {
      height: 1px;
      position: absolute;
      width: 1rem;
      right: -1rem;
      background-color: red;
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
