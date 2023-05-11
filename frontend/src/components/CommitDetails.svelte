<script lang="ts">
  import type { Commit } from "src/scripts/graph";
  import { GetCommitDetails } from "../../wailsjs/go/main/App";

  interface CommitDetails {
  	Body: string;
  	CommitterName: string;
  	CommitterEmail: string;
  	CommitterTimestamp: number;
  	CommitterDatetime: string;
  }

  let commit: Commit;
  let commitDetails: CommitDetails;
  let height = document.documentElement.style.getPropertyValue('--commit-details-height-default');

  (window as any).currentCommitDetails = (): string => {
    return commit?.Hash ?? '';
  }

  (window as any).showCommitDetails = (c: Commit) => {
    commit = c;
    document.documentElement.style.setProperty('--commit-details-height', height);
    GetCommitDetails(c.Hash).then(r => {
      switch (r.Response) {
        case "error":
          (window as any).messageModal(r.Message);
          break;

        case "success":
          commitDetails = r.Commit as CommitDetails;
          break;
      }
    })
  };

  (window as any).hideCommitDetails = () => {
    document.documentElement.style.setProperty('--commit-details-height', '0');
    commit = null;
  }
  (window as any).hideCommitDetails();
</script>

<div class="commit-details">
  {#if commit}
    <div class="commit-details__left">
      <table>
        <tr>
          <th>Hash</th>
          <td>{commit.Hash}</td>
        </tr>
        <tr>
          <th>Parents</th>
          <td>
            {#each commit.Parents as p, i}
              {p}
              {#if i < commit.Parents.length}<br>{/if}
            {/each}
          </td>
        </tr>
        <tr>
          <th>Author</th>
          <td>
            {commit.AuthorName}
            {#if commit.AuthorEmail}
              &lt;<a href="mailto:{commit.AuthorEmail}">{commit.AuthorEmail}</a>&gt;
            {/if}
          </td>
        </tr>
        <tr>
          <th>Authored Date</th>
          <td>{commit.AuthorDatetime}</td>
        </tr>
        <tr>
          <th>Committer</th>
          <td>
            {#if commitDetails}
              {commitDetails.CommitterName}
              {#if commitDetails.CommitterEmail}
                &lt;<a href="mailto:{commitDetails.CommitterEmail}">{commitDetails.CommitterEmail}</a>&gt;
              {/if}
            {/if}
          </td>
        </tr>
        <tr>
          <th>Committed Date</th>
          <td>
            {#if commitDetails}
              {commitDetails.CommitterDatetime}
            {/if}
          </td>
        </tr>
        <tr>
          <th>Message</th>
          <td>
            <div class="message__subject">{commit.Subject}</div>
            {#if commitDetails}
              <div class="message__body">{commitDetails.Body}</div>
            {/if}
          </td>
        </tr>
      </table>
    </div>
    <div class="commit-details__right">
      right
    </div>
  {/if}
</div>
