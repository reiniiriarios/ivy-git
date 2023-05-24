<script lang="ts">
  import octicons from '@primer/octicons';

  import CommitLabels from 'components/commits/CommitLabels.svelte';

  import { NUM_COLORS, UNCOMMITED_HASH } from 'scripts/graph';

  import { HEAD, type Commit, commitSignData } from 'stores/commit-data';
  import { currentCommit } from 'stores/commit-details';
  import { settings } from 'stores/settings';

  export let commit: Commit;
  export let signStatus: string;

  function clearActive() {
    let all = document.getElementsByClassName('commit');
    for (let i = 0; i < all.length; i++) {
      all[i].classList.remove('active');
    }
  }

  function mouseShowDetails(e: MouseEvent & { currentTarget: HTMLElement }) {
    if (e.currentTarget.dataset.uncommitted === 'true') {
      return;
    }
    toggleCommitDetails(e.currentTarget);
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
      toggleCommitDetails(e.currentTarget);
    } else if ($currentCommit.Hash) {
      clearActive();
      currentCommit.unset();
    }
  }

  function toggleCommitDetails(el: HTMLElement) {
    clearActive();
    if ($currentCommit.Hash !== commit.Hash) {
      el.classList.add('active');
    }
    currentCommit.toggle(commit);
  }

  function codify(s: string): string {
    s = s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    return s.replaceAll(/`([^`]+?)`/g, '<code>$1</code>');
  }
</script>

<tr class="commit c-{commit.Color % NUM_COLORS}"
  class:uncommitted={commit.Hash === UNCOMMITED_HASH}
  class:merge={commit.Merge}
  class:stash={commit.Stash}
  class:head={commit.Hash === $HEAD.Hash}
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{commit.Hash === $HEAD.Hash}"
  data-uncommitted="{commit.Hash === UNCOMMITED_HASH}"
  data-menu="{commit.Hash === UNCOMMITED_HASH ? '' : 'commit'}"
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
  <td class="commit__td commit__td--subject" tabindex="0">{@html codify(commit.Subject)}</td>
  {#if $settings.DisplayCommitSignatureInList}
    <td class="commit__td commit__td--gpg">
      {#if commit.Hash === UNCOMMITED_HASH}
        *
      {:else if signStatus && signStatus !== "L" && signStatus !== "N"}
        <span class="gpg-status gpg-status--{signStatus}">{@html octicons.verified.toSVG({width: 16})}</span>
      {/if}
    </td>
  {/if}
  <td class="commit__td commit__td--author">{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td class="commit__td commit__td--authortime">{commit.AuthorDatetime}</td>
</tr>
