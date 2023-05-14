<script lang="ts">
  import CommitDetailsFiles from 'components/CommitDetailsFiles.svelte';
  import { resetDirs } from 'scripts/commit-details-collapse';
  import { resetDetailsSizing, setDetailsResizable } from 'scripts/commit-details-resize';
  import { type Commit } from 'stores/commit-data';
  import { currentCommit, commitDetails, commitDiffSummary } from 'stores/commit-details';

  let height = document.documentElement.style.getPropertyValue('--commit-details-height-default');

  let commit: Commit;
  currentCommit.subscribe(c => {
    resetDirs();
    commit = c;
    if (c?.Hash) {
      document.documentElement.style.setProperty('--commit-details-height', height);
    }
    else {
      document.documentElement.style.setProperty('--commit-details-height', '0');
      resetDetailsSizing();
    }
  });
</script>

<div class="commit-details" use:setDetailsResizable>
  {#if commit?.Hash}
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
            {$commitDetails?.CommitterName}
            {#if $commitDetails?.CommitterEmail}
              &lt;<a href="mailto:{$commitDetails.CommitterEmail}">{$commitDetails.CommitterEmail}</a>&gt;
            {/if}
          </td>
        </tr>
        <tr>
          <th>Committed Date</th>
          <td>
            {$commitDetails?.CommitterDatetime}
          </td>
        </tr>
        <tr>
          <th>Message</th>
          <td>
            <div class="message__subject">{commit.Subject}</div>
            <div class="message__body">{$commitDetails?.Body}</div>
          </td>
        </tr>
      </table>
    </div>
    <div class="commit-details__right">
      <div class="filestatdir">
        <CommitDetailsFiles files={$commitDiffSummary} />
      </div>
    </div>
  {/if}
</div>
