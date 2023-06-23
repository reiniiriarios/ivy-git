<script lang="ts">
  import octicons from '@primer/octicons';
  import CommitLink from 'components/CommitLink.svelte';
  import CommitFileChanges from 'components/commit-details/CommitFileChanges.svelte';
  import { resetDirs } from 'scripts/commit-details-collapse';
  import { resetDetailsSizing, setDetailsResizable } from 'scripts/commit-details-resize';
  import { type Commit } from 'stores/commits';
  import { currentCommit, commitDetails, commitDiffSummary, commitSignature } from 'stores/commit-details';
  import { commitDetailsWindow } from 'stores/ui';
  import SignatureDetails from './SignatureDetails.svelte';
  import type { ComponentConstructorOptions } from 'svelte';
  import CommitMessage from './CommitMessage.svelte';

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
              <CommitLink hash={p} />
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
            <SignatureDetails />
          </td>
        </tr>
        <tr>
          <th>Message</th>
          <td class="commit-details__message">
            <CommitMessage subject={commit.Subject} body={$commitDetails?.BodyHtml} />
          </td>
        </tr>
      </table>
    </div>
    <div class="commit-details__right">
      <div class="filestatdir">
        <CommitFileChanges hash={commit.Hash} files={$commitDiffSummary} />
      </div>
    </div>
  {/if}
</div>