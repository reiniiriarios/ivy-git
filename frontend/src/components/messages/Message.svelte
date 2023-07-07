<script lang="ts">
  import { messageDialog } from 'stores/message-dialog';
  import MessageTag from 'components/messages/MessageTag.svelte';
  import MessageNewRepo from 'components/messages/MessageNewRepo.svelte';
  import MessageNewRemote from 'components/messages/MessageNewRemote.svelte';
  import Checkbox from 'components/elements/Checkbox.svelte';
  import MessageCloneRepo from './MessageCloneRepo.svelte';
  import TextInput from 'components/elements/TextInput.svelte';

  let blankValid: boolean;
  let blankValue: string;

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

  const validateBlank = (value: string): boolean => {
    if ($messageDialog.validateBlank) {
      return !value ? true : $messageDialog.validateBlank(value);
    }
    return true;
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
        {#if $messageDialog.addTag}
          <MessageTag />
        {:else if $messageDialog.newRepo}
          <MessageNewRepo />
        {:else if $messageDialog.cloneRepo}
          <MessageCloneRepo />
        {:else if $messageDialog.addRemote}
          <MessageNewRemote />
        {:else}
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
              <TextInput
                use={(e) => e.focus()}
                display={$messageDialog.blank}
                id="message-dialog-blank"
                validate={validateBlank}
                autoEditBlank={$messageDialog.autoEditBlank}
                bind:value={blankValue}
                bind:valid={blankValid}
              />
            </div>
          {/if}
          {#if $messageDialog.checkboxes}
            <div class="modal__checkboxes">
              {#each $messageDialog.checkboxes as checkbox}
                <Checkbox display={checkbox.label} checked={checkbox.checked} id="checkbox-{checkbox.id}" />
              {/each}
            </div>
          {/if}
          <div class="modal__response">
            {#if $messageDialog.confirm}
              <button
                class="btn yes"
                on:click={messageDialog.yes}
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
