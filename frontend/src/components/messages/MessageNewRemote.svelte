<script lang="ts">
  import { messageDialog } from "stores/message-dialog";

  let remoteName: string;
  let fetchUrl: string;
  let pushUrl: string;
  let nameValid: boolean = false;
  let fetchValid: boolean = false;
  let pushValid: boolean = false;

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  messageDialog.subscribe(() => {
    remoteName = null;
    fetchUrl = null;
    pushUrl = null;
    nameValid = false;
  });

  const validateName = () => nameValid = /^[a-z0-9_\-]+?$/i.test(remoteName);
  const validateFetch = () => fetchValid = validUrl(fetchUrl);
  const validatePush = () => pushValid = validUrl(pushUrl);

  const validUrl = (url: string): boolean => {
    // This is very difficult to validate for all possible URLs
    // https://www.git-scm.com/docs/git-clone#_git_urls
    // Let's just check if it doesn't start or end with a space.
    // Do something with this later.
    return (url[0] != " " && url[url.length-1] != " ");
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if (['\n', 'Enter'].includes(e.key) && remoteName && nameValid && fetchUrl  && fetchValid && (!pushUrl || pushValid)) {
      messageDialog.yes();
    }
  });
</script>

<div class="modal__add-tag">
  <label class="blank-field">
    <span>Remote Name</span>
    <input
      use:focusBlank
      type="text"
      id="message-dialog-remote-name"
      class:invalid={remoteName && !nameValid}
      bind:value={remoteName}
      on:input={validateName}
    >
  </label>
  <label class="blank-field">
    <span>Fetch URL</span>
    <input
      type="text"
      id="message-dialog-remote-fetch"
      class:invalid={fetchUrl && !fetchValid}
      bind:value={fetchUrl}
      on:input={validateFetch}
    >
  </label>
  <label class="blank-field">
    <span>Push URL (Optional)</span>
    <input
      type="text"
      id="message-dialog-remote-push"
      class:invalid={pushUrl && !pushValid}
      bind:value={pushUrl}
      on:input={validatePush}
    >
  </label>
</div>
<div class="modal__response">
  {#if $messageDialog.confirm}
    <button
      class="btn yes"
      on:click={messageDialog.yes}
      disabled={!remoteName || !nameValid || !fetchUrl || !fetchValid || (pushUrl && !pushValid)}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
