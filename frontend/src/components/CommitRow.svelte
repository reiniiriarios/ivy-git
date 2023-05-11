<script lang="ts">
  import { UNCOMMITED_HASH, type Commit, type Ref } from "../scripts/graph";
  import CommitLabels from "./CommitLabels.svelte";

  export let commit: Commit;
  export let HEAD: Ref;

  let u = commit.Hash === UNCOMMITED_HASH;
  let h = commit.Hash === HEAD.Hash;

  let detailsOpen = false;

  function clearActive() {
    let all = document.getElementsByClassName('commit');
    for (let i = 0; i < all.length; i++) {
      all[i].classList.remove('active');
    }
  }

  function showDetails() {
    clearActive();
    if (detailsOpen) {
      (window as any).hideCommitDetails();
    } else {
      (window as any).showCommitDetails(commit);
      (document.querySelector(`.commit[data-id="${commit.Id}"]`) as HTMLElement).classList.add('active');
    }
    detailsOpen = !detailsOpen;
  }
</script>

<tr class="commit c-{commit.Color} {u ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}"
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{h}"
  data-menu="{u ? '' : 'commit'}"
  on:mousedown={showDetails}
  on:keyup={showDetails}>
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
