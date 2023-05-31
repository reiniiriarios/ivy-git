<script lang="ts">
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

  const validateRef = (e: InputEvent & { currentTarget: EventTarget & HTMLInputElement }) => {
    tagValid = checkRef(e.currentTarget.value);
  }

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if (['\n', 'Enter'].includes(e.key) && tagNameField && tagValid && (!annotated || tagMessageField)) {
      messageDialog.yes();
    }
  });
</script>

{#if $messageDialog.heading}
  <div class="modal__heading">{$messageDialog.heading}</div>
{/if}
{#if $messageDialog.message}
  <div class="modal__text">{@html $messageDialog.message}</div>
{/if}
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
  <label class="checkbox">
    <input type="checkbox" id="message-dialog-tag-push">
    <span></span> Push to Remote
  </label>
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
