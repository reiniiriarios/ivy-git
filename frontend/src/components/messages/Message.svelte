<script lang="ts">
  import { messageDialog } from 'stores/message-dialog';
  import MessageTag from 'components/messages/MessageTag.svelte';

  let confirmButton: HTMLButtonElement;

  let blankValue: string;
  let blankValid: boolean = true;

  messageDialog.subscribe(() => {
    blankValue = null;
  });

  window.addEventListener('keydown', function(e: KeyboardEvent) {
    if (['Escape'].includes(e.key) && ($messageDialog.message || $messageDialog.options?.length)) {
      messageDialog.okay();
    }
    else if (['\n', 'Enter'].includes(e.key) && $messageDialog.blank && blankValue.length) {
      messageDialog.yes();
    }
  });

  const focusBlank = (e: HTMLInputElement) => {
    e.focus();
  }

  const validateBlank = (e: InputEvent & { currentTarget: EventTarget & HTMLInputElement }) => {
    if ($messageDialog.validateBlank) {
      blankValid = !e.currentTarget.value ? true : $messageDialog.validateBlank(e.currentTarget.value);
    }
  }
</script>

{#if $messageDialog.message || $messageDialog.options?.length}
  <div role="dialog" class="modal" id="modal-message">
    <div class="overlay">
      <div class="modal__box">
        {#if $messageDialog.addTag}
          <MessageTag />
        {:else}
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
          {/if}
          {#if $messageDialog.blank}
            <div class="modal__blank">
              <label class="blank-field">
                <span>{$messageDialog.blank}</span>
                <input
                  use:focusBlank
                  type="text"
                  id="message-dialog-blank"
                  bind:value={blankValue}
                  class:invalid={blankValue && !blankValid}
                  on:input={validateBlank}
                >
              </label>
            </div>
          {/if}
          {#if $messageDialog.checkboxes}
            <div class="modal__checkboxes">
              {#each $messageDialog.checkboxes as checkbox}
                <label class="checkbox">
                  <input type="checkbox" name="{checkbox.id}" id="checkbox-{checkbox.id}" checked={checkbox.checked}>
                  <span></span>
                  {checkbox.label}
                </label>
              {/each}
            </div>
          {/if}
          <div class="modal__response">
            {#if $messageDialog.confirm}
              <button
                class="btn yes"
                on:click={messageDialog.yes}
                bind:this={confirmButton}
                disabled={$messageDialog.validateBlank ? !blankValue || !blankValid : false}
              >
                {$messageDialog.confirm}
              </button>
            {/if}
            <button class="btn okay" on:click={messageDialog.okay}>
              {$messageDialog.okay}
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
