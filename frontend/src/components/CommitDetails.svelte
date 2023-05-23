<script lang="ts">
  import octicons from '@primer/octicons';
  import CommitDetailsFiles from 'components/CommitDetailsFiles.svelte';
  import { resetDirs } from 'scripts/commit-details-collapse';
  import { resetDetailsSizing, setDetailsResizable } from 'scripts/commit-details-resize';
  import { type Commit } from 'stores/commit-data';
  import { currentCommit, commitDetails, commitDiffSummary, commitSignature } from 'stores/commit-details';
  import { commitDetailsWindow } from 'stores/ui';

  let height = document.documentElement.style.getPropertyValue('--commit-details-height-default');

  let commit: Commit;
  currentCommit.subscribe(c => {
    resetDirs();
    commit = c;
    if (c?.Hash) {
      document.documentElement.style.setProperty('--commit-details-height', height);
      commitDetailsWindow.set(true);
    }
    else {
      document.documentElement.style.setProperty('--commit-details-height', '0px');
      commitDetailsWindow.set(false);
      resetDetailsSizing();
    }
  });

  function codify(s: string): string {
    s = s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    return s.replaceAll(/`([^`]+?)`/g, '<code>$1</code>');
  }
</script>

<div class="commit-details" use:setDetailsResizable class:hidden={!$commitDetailsWindow}>
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
          <th>Signature</th>
          <td class="commit-details__gpg">
            {#if $commitSignature?.Status}
              {#if $commitSignature.Status !== 'N'}
                <span class="gpg-status gpg-status--{$commitSignature.Status}">
                  {@html octicons.verified.toSVG({width: 16})}
                </span>
                <span class="commit-details__gpg-key">{$commitSignature.Key}</span>
                {#if $commitSignature.Name}
                  <span class="commit-details__gpg-name">{$commitSignature.Name}</span>
                {/if}
              {:else}
                None
              {/if}
            {/if}
          </td>
        </tr>
        <tr>
          <th>Message</th>
          <td class="message">
            <div class="message__subject">{@html codify(commit.Subject)}</div>
            <div class="message__body">{@html $commitDetails?.BodyHtml}</div>
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
