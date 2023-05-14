<script lang="ts">
  import { UNCOMMITED_HASH } from 'scripts/graph';
  import CommitLabels from 'components/CommitLabels.svelte';
  import { HEAD, type Commit } from 'stores/commit-data';
  import { currentCommit } from 'stores/commit-details';

  export let commit: Commit;

  let u = commit.Hash === UNCOMMITED_HASH;
  let h = commit.Hash === $HEAD.Hash;

  function clearActive() {
    let all = document.getElementsByClassName('commit');
    for (let i = 0; i < all.length; i++) {
      all[i].classList.remove('active');
    }
  }

  function mouseShowDetails(e: MouseEvent & { currentTarget: HTMLElement }) {
    toggleCommitDetails(e.currentTarget);
  }

  function keyShowDetails(e: KeyboardEvent & { currentTarget: HTMLElement }) {
    if (![' ', 'Enter', 'ArrowUp', 'ArrowRight', 'ArrowDown', 'ArrowLeft'].includes(e.key)) {
      return;
    }
    e.preventDefault();
    if ([' ', 'Enter'].includes(e.key)) {
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
</script>

<tr class="commit c-{commit.Color}"
  class:uncommitted={u}
  class:merge={commit.Merge}
  class:stash={commit.Stash}
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{h}"
  data-menu="{u ? '' : 'commit'}"
  tabindex="0"
  on:click={mouseShowDetails}
  on:keydown={keyShowDetails}>
  <td class="commit__td commit__td--refs">
    {#if commit.Labeled}
      <CommitLabels commit={commit} />
    {/if}
  </td>
  <td class="commit__td commit__td--tree"></td>
  <td class="commit__td commit__td--subject">{commit.Subject}</td>
  <td class="commit__td commit__td--author">{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td class="commit__td commit__td--authortime">{commit.AuthorDatetime}</td>
</tr>
