<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { messageDialog } from "stores/message-dialog";
  import { DirExists, SelectDirectory } from "wailsjs/go/main/App";

  let repoName: string;
  let repoLocation: string;
  let repoValid: boolean = false;

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  function validateName() {
    repoValid = /^[a-z0-9_\-]+?$/i.test(repoName);
    if (repoValid && repoLocation) {
      DirExists(repoName, repoLocation).then(exists => repoValid = !exists);
    }
  }

  function chooseDir() {
    SelectDirectory().then(result => {
      parseResponse(result, () => {
        if (result.Response !== 'none') {
          repoLocation = result.Data;
          if (repoValid && repoName) {
            DirExists(repoName, repoLocation).then(exists => repoValid = !exists);
          }
        }
      });
    });
  }
</script>

{#if $messageDialog.heading}
  <div class="modal__heading">{$messageDialog.heading}</div>
{/if}
{#if $messageDialog.message}
  <div class="modal__text">{@html $messageDialog.message}</div>
{/if}

<div class="modal__new-repo">
  <label class="blank-field">
    <span>Repo Name</span>
    <input
      use:focusBlank
      type="text"
      id="message-dialog-repo-name"
      class:invalid={repoName && !repoValid}
      bind:value={repoName}
      on:input={validateName}
    >
  </label>
  <label class="blank-field">
    <span>Repo Location</span>
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
      disabled={!repoName || !repoValid}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
