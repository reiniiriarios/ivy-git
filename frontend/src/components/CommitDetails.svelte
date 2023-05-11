<script lang="ts">
  import type { Commit } from "src/scripts/graph";

  let commit: Commit;
  let height = document.documentElement.style.getPropertyValue('--commit-details-height-default');

  (window as any).showCommitDetails = (c: Commit) => {
    commit = c;
    document.documentElement.style.setProperty('--commit-details-height', height);
  };

  (window as any).hideCommitDetails = () => {
    document.documentElement.style.setProperty('--commit-details-height', '0');
    commit = null;
  }
  (window as any).hideCommitDetails();
</script>

<div class="commit-details">
  {#if commit}
    <div class="commit-details__left">
      left
      {commit.Hash}
    </div>
    <div class="commit-details__right">
      right
    </div>
  {/if}
</div>
