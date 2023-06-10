<script lang="ts">
  import FileInput from "components/elements/FileInput.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { appData } from "stores/app-data";
  import { messageDialog } from "stores/message-dialog";
  import { DirExists } from "wailsjs/go/main/App";

  let repoName: string;
  let repoLocation: string = $appData.RecentRepoDir;
  let repoValid: boolean = false;

  const validateName = (name: string) => /^[a-z0-9_\-]+?$/i.test(name);

  const validateRepoNameDir = () => {
    DirExists(repoName, repoLocation).then(exists => repoValid = !exists);
  }
</script>

<div class="modal__new-repo">
  <TextInput
    use={(e) => e.focus()}
    display="Repo Name"
    id="message-dialog-repo-name"
    validate={validateName}
    bind:value={repoName}
    bind:valid={repoValid}
  />
  <FileInput
    display="Repo Location"
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
      disabled={!repoName || !repoValid}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
