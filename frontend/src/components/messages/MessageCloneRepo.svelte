<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { appData } from "stores/app-data";
  import { messageDialog } from "stores/message-dialog";
  import { DirExists, SelectDirectory } from "wailsjs/go/main/App";

  let repoUrl: string;
  let repoLocation: string = $appData.RecentRepoDir;
  let repoValid: boolean = false;

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  const validateUrl = () => {
    // This is very difficult to validate for all possible URLs
    // https://www.git-scm.com/docs/git-clone#_git_urls
    // Let's just check if it doesn't start or end with a space.
    // Do something with this later.
    repoValid = (repoUrl[0] != " " && repoUrl[repoUrl.length-1] != " ");
  }

  function chooseDir() {
    SelectDirectory().then(result => {
      parseResponse(result, () => {
        if (result.Response !== 'none') {
          repoLocation = result.Data;
          if (repoValid && repoUrl) {
            let repoName = repoUrl.replace(/^.*\/([^\/]+?)\/?(?:\.git\/?)?$/i, '$1');
            DirExists(repoName, repoLocation).then(exists => repoValid = !exists);
          }
        }
      });
    });
  }
</script>

<div class="modal__clone-repo">
  <label class="blank-field">
    <span>Repo URL</span>
    <input
      use:focusBlank
      type="text"
      id="message-dialog-repo-url"
      class:invalid={repoUrl && !repoValid}
      bind:value={repoUrl}
      on:input={validateUrl}
    >
  </label>
  <label class="blank-field">
    <span>Clone To</span>
    <input
      type="text"
      id="message-dialog-repo-location"
      bind:value={repoLocation}
      on:click={chooseDir}
      readonly
    >
    <button class="btn" id="message-dialog-repo-location-btn" on:click={chooseDir}>Choose</button>
  </label>
</div>

<div class="modal__response">
  {#if $messageDialog.confirm}
    <button
      class="btn yes"
      on:click={messageDialog.yes}
      disabled={!repoUrl || !repoValid}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
