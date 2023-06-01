<script lang="ts">
  import { commitData, commitsPage } from "stores/commit-data";

  function loadMoreCommits() {
    commitData.refresh($commitsPage + 1);
  }

  commitsPage.subscribe(p => {
    let scrollDiv = document.getElementById('commits__scroll');
    if (p && scrollDiv) {
      scrollDiv.scrollTop = scrollDiv.scrollHeight;
    }
  })
</script>

<div class="load-more-commits">
  <button class="btn" on:click={loadMoreCommits}>Load More Commits</button>
</div>
