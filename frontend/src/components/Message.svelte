<script lang="ts">
  import { messageDialog } from 'stores/message-dialog';

  let tagMessageField: HTMLElement;

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if(['Escape'].includes(e.key) && ($messageDialog.message || $messageDialog.options?.length)) {
      messageDialog.okay();
    }
  });

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  const tagAnnotated = () => {
    tagMessageField.style.display = 'block';
  }

  const tagLightweight = () => {
    tagMessageField.style.display = 'none';
  }
</script>

{#if $messageDialog.message || $messageDialog.options?.length}
  <div role="dialog" class="modal" id="modal-message">
    <div class="overlay">
      <div class="modal__box">
        {#if $messageDialog.heading}
          <div class="modal__heading">{$messageDialog.heading}</div>
        {/if}
        {#if $messageDialog.message}
          <div class="modal__text">{@html $messageDialog.message}</div>
        {/if}
        {#if $messageDialog.options}
          <div class="modal__options">
            {#each $messageDialog.options as option}
              <button class="modal__option btn option" on:click={option.callback}>
                {#if option.icon}
                  <div class="modal__option-icon">
                    {@html option.icon}
                  </div>
                {/if}
                {@html option.text}
              </button>
            {/each}
          </div>
        {:else if $messageDialog.checkboxes}
          <div class="modal__checkboxes">
            {#each $messageDialog.checkboxes as checkbox}
              <label class="checkbox">
                <input type="checkbox" name="{checkbox.id}" id="checkbox-{checkbox.id}" checked={checkbox.checked}>
                <span></span>
                {checkbox.label}
              </label>
            {/each}
          </div>
        {:else if $messageDialog.blank}
          <div class="modal__blank">
            <label class="blank-field">
              <span>{$messageDialog.blank}</span>
              <input use:focusBlank type="text" id="message-dialog-blank">
            </label>
          </div>
        {:else if $messageDialog.addTag}
          <div class="modal__add-tag">
            <label class="blank-field">
              <span>Tag Name</span>
              <input use:focusBlank type="text" id="message-dialog-tag-name">
            </label>
            <div class="radio">
              <span class="radio__label">Type</span>
              <label class="radio__option">
                <input type="radio" value="annotated" name="message-dialog-tag-type" checked on:click={tagAnnotated}><span></span> Annotated
              </label>
              <label class="radio__option">
                <input type="radio" value="lightweight" name="message-dialog-tag-type" on:click={tagLightweight}><span></span> Lightweight
              </label>
            </div>
            <label class="blank-field" bind:this={tagMessageField}>
              <span>Message</span>
              <input type="text" id="message-dialog-tag-message">
            </label>
            <label class="checkbox">
              <input type="checkbox" id="message-dialog-tag-push">
              <span></span> Push to Remote
            </label>
          </div>
        {/if}
        <div class="modal__response">
          {#if $messageDialog.confirm}
            <button class="btn yes" on:click={messageDialog.yes}>{$messageDialog.confirm}</button>
          {/if}
          <button class="btn okay" on:click={messageDialog.okay}>{$messageDialog.okay}</button>
        </div>
      </div>
    </div>
  </div>
{/if}
