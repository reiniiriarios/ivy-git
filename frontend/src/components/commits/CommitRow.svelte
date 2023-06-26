<script lang="ts">
  import octicons from '@primer/octicons';

  import CommitLabels from 'components/commits/CommitLabels.svelte';
  import { highlightConventionalCommits } from 'scripts/conventional-commits';

  import { NUM_COLORS, UNCOMMITED_HASH } from 'scripts/graph';

  import { HEAD, type Commit } from 'stores/commits';
  import { currentCommit } from 'stores/commit-details';
  import { RepoState, repoState } from 'stores/repo-state';
  import { settings } from 'stores/settings';
  import { avatars } from 'stores/avatars';

  export let commit: Commit;
  export let signStatus: string;

  function mouseShowDetails(e: MouseEvent & { currentTarget: HTMLElement }) {
    if (e.currentTarget.dataset.uncommitted === 'true') {
      return;
    }
    currentCommit.toggle(commit);
  }

  function keyShowDetails(e: KeyboardEvent & { currentTarget: HTMLElement }) {
    if (e.currentTarget.dataset.uncommitted === 'true') {
      return;
    }
    if (![' ', '\n', 'Enter', 'ArrowUp', 'ArrowRight', 'ArrowDown', 'ArrowLeft'].includes(e.key)) {
      return;
    }
    e.preventDefault();
    if ([' ', '\n', 'Enter'].includes(e.key)) {
      currentCommit.toggle(commit);
    } else if ($currentCommit.Hash) {
      currentCommit.clear();
    }
  }

  function formatSubject(s: string): string {
    s = codify(s);
    if ($settings.HighlightConventionalCommits) {
      s = highlightConventionalCommits(s);
    }
    return s;
  }

  function codify(s: string): string {
    s = s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    return s.replaceAll(/`([^`]+?)`/g, '<code>$1</code>');
  }

  let avatarUrl: string = '';
  $: {
    // Empty first, then wait on load.
    avatarUrl = '';
    if ($settings.DisplayAvatars && commit.AuthorEmail) {
      avatars.fetch(commit.AuthorEmail).then(url => avatarUrl = url);
    }
  }
</script>

<tr class="commit c-{commit.Color % NUM_COLORS} {commit.Hash === UNCOMMITED_HASH ? `repo-state--${$repoState}` : ''}"
  class:uncommitted={commit.Hash === UNCOMMITED_HASH}
  class:merge={commit.Merge}
  class:stash={commit.Stash}
  class:head={commit.Hash === $HEAD.Hash && [RepoState.Nil, RepoState.None].includes($repoState)}
  class:active={commit.Hash === $currentCommit.Hash}
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{commit.Hash === $HEAD.Hash}"
  data-uncommitted="{commit.Hash === UNCOMMITED_HASH}"
  data-menu="{commit.Hash === UNCOMMITED_HASH ? '' : 'commit'}"
  data-merge="{commit.Merge}"
  on:click={mouseShowDetails}
  on:keydown={keyShowDetails}>
  <td class="commit__td commit__td--refs">
    {#if commit.Labeled}
      <CommitLabels commit={commit} />
    {/if}
  </td>
  <td class="commit__td commit__td--tree"></td>
  <!-- The following is interactive via the tr, which doesn't take a tabindex because it's set to display: contents. -->
  <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
  <td class="commit__td commit__td--subject" tabindex="0">{@html formatSubject(commit.Subject)}</td>
  {#if $settings.DisplayCommitSignatureInList}
    <td class="commit__td commit__td--gpg">
      {#if commit.Hash === UNCOMMITED_HASH}
        *
      {:else if signStatus && signStatus !== "L" && signStatus !== "N"}
        <span class="gpg-status gpg-status--{signStatus}">{@html octicons.verified.toSVG({width: 16})}</span>
      {/if}
    </td>
  {/if}
  <td class="commit__td commit__td--author">
    {#if avatarUrl}
      <span class="avatar"><img src="{avatarUrl}" alt="" /></span>
    {/if}
    {commit.AuthorName ?? commit.AuthorEmail}
  </td>
  <td class="commit__td commit__td--authortime">{commit.AuthorDatetime}</td>
</tr>
