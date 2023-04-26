<script lang="ts">
  import {
    GetRepos,
    AddRepo,
    UpdateSelectedRepo,
    GetSelectedRepo,
    RemoveRepo,
  } from "../../wailsjs/go/main/App";

  interface Repo {
    Name: string;
    Directory: string;
  }

  let selectedRepo: string;
  let repos: Repo[] = [];
  let listVisible: boolean = false;

  function getRepos() {
    GetRepos().then((result) => (repos = result as Repo[]));
  }
  getRepos();

  function getSelectedRepo() {
    GetSelectedRepo().then((result) => (selectedRepo = result));
  }
  getSelectedRepo();

  function addRepo() {
    AddRepo().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          repos[result.Repo.Id as string] = result.Repo as Repo;
          repos = repos;
          break;

        case "":
        default:
          // nothing, no dir was selected
      }
    });
  }

  function selectRepo(e: any) {
    UpdateSelectedRepo(e.target.dataset.id).then(() => {
      selectedRepo = e.target.dataset.id;
      hideList();
    });
  }

  function toggleList() {
    if (listVisible) {
      hideList();
    } else {
      showList();
    }
  }

  function showList() {
    document.getElementById("all-repos").style.display = "block";
    listVisible = true;
  }

  function hideList() {
    document.getElementById("all-repos").style.display = "none";
    listVisible = false;
  }

  function delRepo(e: any) {
    var name = repos[e.target.dataset.id].Name;
    (window as any).confirmModal(`Are you sure you want to remove ${name}?`, () => {
      RemoveRepo(e.target.dataset.id).then((result) => (repos = result as Repo[]));
    }, 'Remove', 'Cancel');
  }
</script>

<button class="btn" id="current-repo" on:click={toggleList} on:keyup={toggleList}>
  <div class="label">Current Repo:</div>
  <div>{repos[selectedRepo] ? repos[selectedRepo].Name : 'none selected'}</div>
</button>

<div id="all-repos">
  <div class="overlay" on:click={hideList} on:keyup={hideList}></div>
  <div id="all-repos-bar">
    <div id="add-repo-bar">
      <button class="btn" id="add-repo" on:click={addRepo} on:keyup={addRepo}>Add Repo +</button>
    </div>
    <ul id="repos-list">
      {#each Object.entries(repos) as [id, repo]}
        <li>
          <button class="name" on:click={selectRepo} data-id={id}>{repo.Name}</button>
          <button class="x" on:click={delRepo} on:keyup={delRepo} data-id={id}>&times;</button>
        </li>
      {/each}
    </ul>
  </div>
</div>

<style lang="scss">
  #current-repo {
    height: 4rem;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    text-align: left;
    padding-left: 1.5rem;

    .label {
      color: var(--color-text-label);
      font-size: 0.8rem;
    }
  }

  #all-repos {
    display: none;

    .overlay {
      left: 20rem;
      width: calc(100% - 20rem);
    }

    #all-repos-bar {
      background-color: var(--color-sidebar-bg);
      position: absolute;
      top: 4rem;
      left: 0;
      height: calc(100% - 4rem);
      width: var(--sidebar-width);
    }
  }

  #add-repo {
    width: 100%;
    border-top: 1px solid var(--color-btn-border);
  }

  #repos-list {
    list-style: none;
    margin: 0;
    padding: 0;

    li {
      margin: 0;
      padding: 0;
      display: flex;
      justify-content: space-between;

      button {
        text-align: left;
        border: 0;
        font-size: 0.9rem;
        padding: 0.6rem 0.75rem;
        margin: 0;
        background-color: transparent;
        color: var(--color-text);
        cursor: pointer;

        &:hover {
          background-color: var(--color-btn-bg-hover);
        }

        &.name {
          width: 100%;
        }

        &.x {
          &:hover {
            background-color: var(--color-red-800);
          }
        }
      }
    }
  }
</style>
