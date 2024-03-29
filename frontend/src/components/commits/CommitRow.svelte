<script lang="ts">
  import octicons from '@primer/octicons';

  import CommitLabels from 'components/commits/CommitLabels.svelte';
  import Avatar from 'components/elements/Avatar.svelte';

  import { NUM_COLORS, UNCOMMITED_HASH } from 'scripts/graph';
  import { highlightConventionalCommits } from 'scripts/conventional-commits';

  import { HEAD, type Commit } from 'stores/commits';
  import { currentCommit } from 'stores/commit-details';
  import { RepoState, repoState } from 'stores/repo-state';
  import { settings } from 'stores/settings';
  import { uncommittedChanges } from 'stores/changes';

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
</script>

<tr class="commit c-{commit.Color % NUM_COLORS} {commit.Hash === UNCOMMITED_HASH ? `repo-state--${$repoState}` : ''}"
  class:uncommitted={commit.Hash === UNCOMMITED_HASH}
  class:merge={commit.Merge}
  class:stash={commit.Stash}
  class:head={commit.Hash === $HEAD.Hash && !$uncommittedChanges && [RepoState.Nil, RepoState.None].includes($repoState)}
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
  <td class="commit__td commit__td--subject commit__ellipsis" tabindex="0">
    {@html formatSubject(commit.Subject)}
  </td>
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
    {#if $settings.DisplayAvatars && commit.AuthorEmail}
      <Avatar email="{commit.AuthorEmail}" />
    {/if}
    <span class="commit__ellipsis">
      {commit.AuthorName ?? commit.AuthorEmail}
    </span>
  </td>
  <td class="commit__td commit__td--authortime commit__ellipsis">
    {commit.AuthorDatetime}
  </td>
</tr>
