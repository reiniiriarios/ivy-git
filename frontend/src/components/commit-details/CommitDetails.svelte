<script lang="ts">
  import CommitLink from 'components/elements/CommitLink.svelte';
  import CommitFileChanges from 'components/commit-details/CommitFileChanges.svelte';
  import { resetDirs } from 'scripts/commit-details-collapse';
  import { resetDetailsSizing, setDetailsResizable } from 'scripts/commit-details-resize';
  import { type Commit } from 'stores/commits';
  import { currentCommit, commitDetails, commitDiffSummary } from 'stores/commit-details';
  import { commitDetailsWindow } from 'stores/ui';
  import SignatureDetails from 'components/commit-details/SignatureDetails.svelte';
  import CommitMessage from 'components/commit-details/CommitMessage.svelte';
  import Avatar from 'components/elements/Avatar.svelte';
  import { settings } from 'stores/settings';
  import { onDestroy } from 'svelte';

  let height = document.documentElement.style.getPropertyValue('--commit-details-height-default');

  let commit: Commit;
  const currentCommitUnsubscribe = currentCommit.subscribe(c => {
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

  onDestroy(() => {
    currentCommitUnsubscribe();
  });
</script>

<div class="commit-details" use:setDetailsResizable class:hidden={!$commitDetailsWindow}>
  {#if commit?.Hash}
    <div class="commit-details__left">
      <table>
        <tr>
          <th>Hash</th>
          <td>
            <CommitLink hash={commit.Hash} disabled={true} />
          </td>
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
          <td data-menu="text">
            {#if $settings.DisplayAvatars && commit.AuthorEmail}
              <Avatar email="{commit.AuthorEmail}" hover={true} />
            {/if}
            {commit.AuthorName}
            {#if commit.AuthorEmail}
              &lt;<a href="mailto:{commit.AuthorEmail}">{commit.AuthorEmail}</a>&gt;
            {/if}
          </td>
        </tr>
        <tr>
          <th>Authored Date</th>
          <td data-menu="text">{commit.AuthorDatetime}</td>
        </tr>
        <tr>
          <th>Committer</th>
          <td data-menu="text">
            {#if $settings.DisplayAvatars && $commitDetails?.CommitterEmail}
              <Avatar email="{$commitDetails.CommitterEmail}" hover={true} />
            {/if}
            {$commitDetails?.CommitterName}
            {#if $commitDetails?.CommitterEmail}
              &lt;<a href="mailto:{$commitDetails.CommitterEmail}">{$commitDetails.CommitterEmail}</a>&gt;
            {/if}
          </td>
        </tr>
        <tr>
          <th>Committed Date</th>
          <td data-menu="text">
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
