<script lang="ts">
  import { RepoState, repoState, repoStateMessage } from "stores/repo-state";
</script>

{#if $repoState != RepoState.Nil && $repoState != RepoState.None}
  <div class="repo-state repo-state--{$repoState}">
    <div class="repo-state__banner">
      {$repoStateMessage}
    </div>
  </div>
{/if}
