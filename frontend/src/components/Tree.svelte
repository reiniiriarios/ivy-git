<script lang="ts">
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

  interface Ref {
    Hash: string;
    Name: string;
  }

  let commits: Commit[] = [];
  let HEAD: string = "";

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
            <td class="tree__refs">
              {#if commit.Heads && commit.Heads.length}
                {#each commit.Heads as h}
                  <div class="tree__head">
                    {h}
                    {#if commit.Remotes && commit.Remotes.length}
                      {#each commit.Remotes as r}
                        <div class="tree__remote">{r}</div>
                      {/each}
                    {/if}
                    {#if commit.Hash == HEAD}
                      <div class="tree__HEAD">HEAD</div>
                    {/if}
                  </div>
                {/each}
              {/if}
              {#if commit.Tags && commit.Tags.length}
                {#each commit.Tags as t}
                  <div class="tree__tag">{t}</div>
                {/each}
              {/if}
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

<style lang="scss">
  .tree {
    width: 100%;
    height: 100%;

    &__table {
      width: 100%;

      th, td {
        text-align: left;
        padding: 0.125rem 0.5rem;
      }
    }
  }
</style>
