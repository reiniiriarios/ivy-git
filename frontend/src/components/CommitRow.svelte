<script lang="ts">
  import { UNCOMMITED_HASH, type Commit, type Ref } from "../scripts/graph";
  import CommitLabels from "./CommitLabels.svelte";

  export let commit: Commit;
  export let HEAD: Ref;

  let u = commit.Hash === UNCOMMITED_HASH;
  let h = commit.Hash === HEAD.Hash;

  function clearActive() {
    let all = document.getElementsByClassName('commit');
    for (let i = 0; i < all.length; i++) {
      all[i].classList.remove('active');
    }
  }

  function showDetails(e: MouseEvent | KeyboardEvent) {
    if (e instanceof KeyboardEvent && ![' ', 'Enter'].includes(e.key)) {
      return;
    }
    clearActive();
    if ((window as any).currentCommitDetails() === commit.Hash) {
      (window as any).hideCommitDetails();
    } else {
      (window as any).showCommitDetails(commit);
      (document.querySelector(`.commit[data-id="${commit.Id}"]`) as HTMLElement).classList.add('active');
    }
  }

  function setKeyboard(el: HTMLElement) {
    window.addEventListener('keydown', function(e: KeyboardEvent) {
      if(e.key === ' ' && e.target === el) {
        e.preventDefault();
      }
    });
  }
</script>

<tr class="commit c-{commit.Color} {u ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}"
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{h}"
  data-menu="{u ? '' : 'commit'}"
  tabindex="0"
  on:click={showDetails}
  on:keyup={showDetails}
  use:setKeyboard>
  <td class="commit__td commit__td--refs">
    {#if commit.Labeled}
      <CommitLabels commit={commit} HEAD={HEAD} />
    {/if}
  </td>
  <td class="commit__td commit__td--tree"></td>
  <td class="commit__td commit__td--subject">{commit.Subject}</td>
  <td class="commit__td commit__td--author">{commit.AuthorName ?? commit.AuthorEmail}</td>
  <td class="commit__td commit__td--authortime">{commit.AuthorDatetime}</td>
</tr>
