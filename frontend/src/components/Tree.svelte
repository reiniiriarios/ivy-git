<script lang="ts">
  import { GetCommitsForTree } from "../../wailsjs/go/main/App";

  interface Commit {
    Hash: string;
    Parent: string;
    AuthorName: string;
    AuthorEmail: string;
    AuthorDate: string; // todo
    Subject: string;
  }

  let commits: Commit[] = [];

  (window as any).GetCommitsForTree = () => {
    GetCommitsForTree().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          commits = result.Commits as Commit[];
          console.log(commits[0]);
          break;
      }
    });
  };
</script>

<div class="tree">
  <div class="tree__graph"></div>
  <div class="tree__table">
    {#if commits.length}
      <table>
        <tr>
          <th>Commit</th>
          <th>Date</th>
        </tr>
        {#each commits as commit}
          <tr>
            <td>{commit.Subject}</td>
            <td>{commit.AuthorDate}</td>
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
  }

  table {
    width: 100%;
  }

  th, td {
    text-align: left;
    padding: 0.125rem 0.5rem;
  }
</style>
