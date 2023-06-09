<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import { checkRef } from "scripts/check-ref";
  import { messageDialog } from "stores/message-dialog";

  let tagMessage: HTMLElement;
  let tagAnnotatedField: HTMLInputElement;
  let tagNameField: string;
  let tagMessageField: string;
  let tagValid: boolean = false;

  let annotated = true;

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  messageDialog.subscribe(() => {
    tagNameField = null;
    tagMessageField = null;
  });

  const tagAnnotated = () => {
    tagMessage.style.display = 'block';
    annotated = true;
  }

  const tagLightweight = () => {
    tagMessage.style.display = 'none';
    annotated = false;
  }

  const validateRef = () => {
    tagValid = checkRef(tagNameField);
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if (['\n', 'Enter'].includes(e.key) && tagNameField && tagValid && (!annotated || tagMessageField)) {
      messageDialog.yes();
    }
  });
</script>

<div class="modal__add-tag">
  <label class="blank-field">
    <span>Tag Name</span>
    <input
      use:focusBlank
      type="text"
      id="message-dialog-tag-name"
      class:invalid={tagNameField && !tagValid}
      bind:value={tagNameField}
      on:input={validateRef}
    >
  </label>
  <div class="radio">
    <span class="radio__label">Type</span>
    <label class="radio__option">
      <input
        type="radio"
        value="annotated"
        name="message-dialog-tag-type"
        checked
        on:click={tagAnnotated}
        bind:this={tagAnnotatedField}
      ><span></span> Annotated
    </label>
    <label class="radio__option">
      <input
        type="radio"
        value="lightweight"
        name="message-dialog-tag-type"
        on:click={tagLightweight}
      ><span></span> Lightweight
    </label>
  </div>
  <label class="blank-field" bind:this={tagMessage}>
    <span>Message</span>
    <input type="text" id="message-dialog-tag-message" bind:value={tagMessageField}>
  </label>
  <Checkbox display="Push to Remote" id="message-dialog-tag-push" />
</div>
<div class="modal__response">
  {#if $messageDialog.confirm}
    <button
      class="btn yes"
      on:click={messageDialog.yes}
      disabled={!tagNameField || !tagValid || (annotated && !tagMessageField)}
    >
      {$messageDialog.confirm}
    </button>
  {/if}
  <button class="btn okay" on:click={messageDialog.okay}>
    {$messageDialog.okay}
  </button>
</div>
