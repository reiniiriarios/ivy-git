<script lang="ts">
  import { UNCOMMITED_HASH, type Commit, type Ref } from "../scripts/graph";
  import CommitLabels from "./CommitLabels.svelte";

  export let commit: Commit;
  export let HEAD: Ref;

  let u = commit.Hash === UNCOMMITED_HASH;
  let h = commit.Hash === HEAD.Hash;
</script>

<tr class="commit c-{commit.Color} {u ? 'uncommitted' : ''} {commit.Merge ? 'merge' : ''} {commit.Stash ? 'stash' : ''}"
  data-id="{commit.Id}"
  data-hash="{commit.Hash}"
  data-head="{h}"
  data-menu="{u ? '' : 'commit'}">
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
