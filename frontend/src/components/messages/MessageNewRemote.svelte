<script lang="ts">
  import TextInput from "components/elements/TextInput.svelte";
  import { messageDialog } from "stores/message-dialog";

  let remoteName: string;
  let fetchUrl: string;
  let pushUrl: string;
  let nameValid: boolean = false;
  let fetchValid: boolean = false;
  let pushValid: boolean = false;

  messageDialog.subscribe(() => {
    remoteName = null;
    fetchUrl = null;
    pushUrl = null;
    nameValid = false;
  });

  const validateName = () => nameValid = /^[a-z0-9_\-]+?$/i.test(remoteName);

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
  <TextInput
    use={(e) => e.focus()}
    display="Remote Name"
    id="message-dialog-remote-name"
    validate={validateName}
    bind:value={remoteName}
    bind:valid={nameValid}
  />
  <TextInput
    display="Fetch URL"
    id="message-dialog-remote-fetch"
    validate={validUrl}
    bind:value={fetchUrl}
    bind:valid={fetchValid}
  />
  <TextInput
    display="Push URL (Optional)"
    id="message-dialog-remote-push"
    validate={validUrl}
    bind:value={pushUrl}
    bind:valid={pushValid}
  />
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
