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

  (window as any).getSelectedRepo = () => {
    GetSelectedRepo().then((r) => {
      selectedRepo = r;
      (window as any).selectedRepo = r;
    });
  }

  (window as any).getRepos = () => {
    GetRepos().then((result) => (repos = result as Repo[]));
  }

  function addRepo() {
    AddRepo().then((result) => {
      switch (result.Response) {
        case "error":
          (window as any).messageModal(result.Message);
          break;

        case "success":
          repos[result.Id as string] = result.Repo as Repo;
          repos = repos;
          break;
      }
    });
  }

  function selectRepo(e: any) {
    UpdateSelectedRepo(e.target.dataset.id).then(() => {
      selectedRepo = e.target.dataset.id;
      (window as any).selectedRepo = e.target.dataset.id;
      (window as any).getCurrentBranch();
      (window as any).getBranches();
      (window as any).getChanges();
      hideList();
      if ((window as any).currentTab == 'tree') {
        (window as any).GetCommitList();
      }
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
    document.getElementById("current-repo").classList.add('active');
    listVisible = true;
  }

  function hideList() {
    document.getElementById("all-repos").style.display = "none";
    document.getElementById("current-repo").classList.remove('active');
    listVisible = false;
  }

  function delRepo(e: any) {
    var name = repos[e.target.dataset.id].Name;
    (window as any).confirmModal(`Are you sure you want to remove ${name}?`, () => {
      RemoveRepo(e.target.dataset.id).then((result) => (repos = result as Repo[]));
    }, 'Remove', 'Cancel');
  }
</script>

<button class="btn btn-drop sidebar-big-button" id="current-repo" on:click={toggleList} on:keyup={toggleList}>
  <div class="sidebar-big-button__label">Current Repo:</div>
  <div class="sidebar-big-button__value">{repos[selectedRepo] ? repos[selectedRepo].Name : 'none selected'}</div>
</button>

<div id="all-repos" class="sidebar-dropdown">
  <div class="overlay" on:click={hideList} on:keyup={hideList}></div>
  <div class="sidebar-dropdown__container">
    <div class="sidebar-dropdown__bar">
      <div class="sidebar-dropdown__add">
        <button class="btn" on:click={addRepo} on:keyup={addRepo}>Add Repo +</button>
      </div>
      <ul class="sidebar-dropdown__list">
        {#each Object.entries(repos) as [id, repo]}
          <li>
            <button class="list-btn name" on:click={selectRepo} data-id={id}>{repo.Name}</button>
            <button class="list-btn x" on:click={delRepo} on:keyup={delRepo} data-id={id}>&times;</button>
          </li>
        {/each}
      </ul>
    </div>
  </div>
</div>
