<script lang="ts">
  import { commitData, commitsPage } from "stores/commits";
  import { onDestroy } from "svelte";

  function loadMoreCommits() {
    commitData.refresh($commitsPage + 1);
  }

  const commitsPageUnsubscribe = commitsPage.subscribe(p => {
    let scrollDiv = document.getElementById('commits__scroll');
    if (p && scrollDiv) {
      scrollDiv.scrollTop = scrollDiv.scrollHeight;
    }
  });

  onDestroy(() => {
    commitsPageUnsubscribe();
  });
</script>

<div class="load-more-commits">
  <button class="btn" on:click={loadMoreCommits}>Load More Commits</button>
</div>
