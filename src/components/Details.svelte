<script lang="ts">
  import Remotes from "./Remotes.svelte";
  import RepoInfo from "./RepoInfo.svelte";
  import Stashes from "./Stashes.svelte";
</script>

<div class="details">
  <div class="details__left">
    <RepoInfo />
  </div>
  <div class="details__right">
    <div class="details__right-section">
      <Remotes />
    </div>
    <div class="details__right-section">
      <Stashes />
    </div>
  </div>
</div>

<style lang="scss">
  .details {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: row;

    &__left,
    &__right {
      flex: 1;
    }

    &__right {
      display: flex;
      flex-direction: column;
    }

    &__right-section {
      flex: 1;
    }
  }
</style>
