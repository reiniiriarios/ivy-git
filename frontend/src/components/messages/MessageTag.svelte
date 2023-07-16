<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { checkRef } from "scripts/check-ref";
  import { messageDialog } from "stores/message-dialog";
  import { onDestroy } from "svelte";

  let tagName: string;
  let tagMessage: string;
  let tagValid: boolean;

  const messageDialogUnsubscribe = messageDialog.subscribe(() => {
    tagName = null;
    tagMessage = null;
  });

  onDestroy(() => {
    messageDialogUnsubscribe();
  });

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if (['\n', 'Enter'].includes(e.key) && tagName && tagValid) {
      messageDialog.yes();
    }
  });
</script>

<div class="modal__add-tag">
  <TextInput
    use={(e) => e.focus()}
    display="Tag Name"
    id="message-dialog-tag-name"
    validate={checkRef}
    bind:value={tagName}
    bind:valid={tagValid}
  />
  <TextInput
    display="Message"
    id="message-dialog-tag-message"
    bind:value={tagMessage}
  />
  <Checkbox display="Push to Remote" id="message-dialog-tag-push" />
</div>
<div class="modal__response">
  {#if $messageDialog.confirm}
    <button
      class="btn yes"
      on:click={messageDialog.yes}
      disabled={!tagName || !tagValid}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
