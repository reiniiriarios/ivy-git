<script lang="ts">
  import FileInput from "components/elements/FileInput.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { appData } from "stores/app-data";
  import { messageDialog } from "stores/message-dialog";
  import { FileExists } from "wailsjs/go/ivy/App";

  let repoUrl: string;
  let repoLocation: string = $appData.RecentRepoDir;
  let repoValid: boolean = false;

  const validateUrl = (value: string) => {
    // This is very difficult to validate for all possible URLs
    // https://www.git-scm.com/docs/git-clone#_git_urls
    // Let's just check if it doesn't start or end with a space.
    // Do something with this later.
    if (value && value[0] != " " && value[value.length-1] != " ") {
      validateRepoNameDir();
    }
    return false;
  }

  const validateRepoNameDir = () => {
    let repoName = repoUrl.replace(/^.*\/([^\/]+?)\/?(?:\.git\/?)?$/i, '$1');
    FileExists(repoName, repoLocation).then(exists => repoValid = !exists);
  }
</script>

<div class="modal__clone-repo">
  <TextInput
    use={(e) => e.focus()}
    display="Repo URL"
    id="message-dialog-repo-url"
    validate={validateUrl}
    bind:value={repoUrl}
    bind:valid={repoValid}
  />
  <FileInput
    display="Clone To"
    id="message-dialog-repo-location"
    directory={true}
    bind:value={repoLocation}
    on:change={validateRepoNameDir}
  />
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
