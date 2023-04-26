<script lang="ts">
  import {
    GetCurrentBranch,
    GetBranches,
  } from "../../wailsjs/go/main/App";

  interface Branch {
    Name: string
  }

  let selectedBranch: Branch;
  let branches: Branch[] = [];
  let listVisible: boolean = false;

  (window as any).getBranches = () => {
    GetBranches().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          branches = result.Branches;
          break;
      }
    });
  }

  (window as any).getCurrentBranch = () => {
    GetCurrentBranch().then((result) => {
      selectedBranch = result.Branch as Branch;
      selectedBranch = selectedBranch;
    });
  }

  function newBranch() {
    // ...
  }

  function switchBranch(e: any) {
    // SelectBranch(e.target.dataset.name).then(() => {
    //   selectedBranch = e.target.dataset.name;
    //   hideList();
    // });
  }

  function toggleList() {
    if (listVisible) {
      hideList();
    } else {
      showList();
    }
  }

  function showList() {
    document.getElementById("all-branches").style.display = "block";
    listVisible = true;
  }

  function hideList() {
    document.getElementById("all-branches").style.display = "none";
    listVisible = false;
  }
</script>

<button class="btn" id="current-branch" on:click={toggleList} on:keyup={toggleList}>
  <div class="label">Current Branch:</div>
  <div>{selectedBranch?.Name ?? 'none selected'}</div>
</button>

<div id="all-branches">
  <div class="overlay" on:click={hideList} on:keyup={hideList}></div>
  <div id="all-branches__container">
    <div id="all-branches__bar">
      <div id="all-branches__add">
        <button class="btn" on:click={newBranch} on:keyup={newBranch}>Create Branch +</button>
      </div>
      <ul id="all-branches__list">
        {#each Object.entries(branches) as [_, branch]}
          <li>
            <button class="list-btn name" on:click={switchBranch} data-name={branch?.Name}>{branch?.Name}</button>
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>

<style lang="scss">
  #current-branch {
    height: 4rem;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    text-align: left;
    padding-left: 1.5rem;
    border-top: 1px solid var(--color-btn-border);

    .label {
      color: var(--color-text-label);
      font-size: 0.8rem;
    }
  }

  #all-branches {
    display: none;
    height: 100%;

    .overlay {
      left: 20rem;
      width: calc(100% - 20rem);
    }

    &__container {
      position: relative;
      height: 100%;
    }

    &__bar {
      background-color: var(--color-sidebar-bg);
      position: absolute;
      top: 0;
      left: 0;
      height: 100%;
      width: var(--sidebar-width);
    }

    &__add {
      width: 100%;
      border-top: 1px solid var(--color-btn-border);

      button {
        width: 100%;
      }
    }

    &__list {
      list-style: none;
      margin: 0;
      padding: 0;

      li {
        margin: 0;
        padding: 0;
        display: flex;
        justify-content: space-between;
      }
    }
  }
</style>
