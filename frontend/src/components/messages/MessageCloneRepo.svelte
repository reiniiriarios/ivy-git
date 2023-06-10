<script lang="ts">
  import TextInput from "components/elements/TextInput.svelte";
  import { parseResponse } from "scripts/parse-response";
  import { appData } from "stores/app-data";
  import { messageDialog } from "stores/message-dialog";
  import { DirExists, SelectDirectory } from "wailsjs/go/main/App";

  let repoUrl: string;
  let repoLocation: string = $appData.RecentRepoDir;
  let repoValid: boolean = false;

  const validateUrl = (value: string) => {
    // This is very difficult to validate for all possible URLs
    // https://www.git-scm.com/docs/git-clone#_git_urls
    // Let's just check if it doesn't start or end with a space.
    // Do something with this later.
    return (value[0] != " " && value[value.length-1] != " ");
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
  <TextInput
    use={(e) => e.focus()}
    display="Repo URL"
    classes="blank-field"
    id="message-dialog-repo-url"
    validate={validateUrl}
    bind:value={repoUrl}
    bind:valid={repoValid}
  />
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
