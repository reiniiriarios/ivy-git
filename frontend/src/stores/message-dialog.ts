import { writable } from 'svelte/store';

type MessageCallback = (e?: MouseEvent & { currentTarget: HTMLElement } | KeyboardEvent & { currentTarget: HTMLElement }) => any;

interface Message {
  heading?: string;
  message?: string;
  confirm?: string;
  callbackConfirm?: MessageCallback;
  okay?: string;
  callback?: MessageCallback;
  options?: {
    text: string;
    icon?: string;
    callback: MessageCallback;
  }[],
  checkboxes?: {
    id: string;
    label: string;
    checked?: boolean;
  }[],
  blank?: string,
}

function createMessage() {
  const { subscribe, set, update } = writable({} as Message);

  return {
    subscribe,
    show: async (message: Message) => {
      set({
        heading: message.heading ?? 'Notice',
        message: message.message ?? '...',
        callback: message.callback ?? (() => {}),
        okay: message.okay ?? 'Okay'
      });
    },
    error: async (message: Message) => {
      set({
        heading: message.heading ?? 'Error',
        message: message.message ?? 'Unknown error occurred.',
        callback: message.callback ?? (() => {}),
        okay: message.okay ?? 'Okay'
      });
    },
    options: async (message: Message) => {
      if (message.options?.length) {
        set({
          heading: message.heading ?? 'Select Option',
          message: message.message ?? '',
          callback: message.callback ?? (() => {}),
          okay: message.okay ?? 'Cancel',
          options: message.options,
        });
      } else {
        console.error('No options given to message dialog.');
      }
    },
    confirm: async(message: Message) => {
      set({
        heading: message.heading ?? 'Confirm',
        message: message.message ?? 'Confirm?',
        confirm: message.confirm ?? 'Yes',
        checkboxes: message.checkboxes ?? [],
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'No',
        callback: message.callback ?? (() => {}),
      });
    },
    fillBlank: async(message: Message) => {
      set({
        heading: message.heading ?? 'Enter Data',
        message: message.message ?? 'Enter data:',
        confirm: message.confirm ?? 'Send',
        blank: message.blank ?? 'Data',
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    yes: async() => {
      update(message => {
        message.callbackConfirm();
        return {};
      })
    },
    okay: async() => {
      update(message => {
        message.callback();
        return {};
      })
    },
    clear: async() => {
      set({});
    },
    // Shortcut for getting form option value.
    tickboxTicked: (id: string) => {
      let el = document.getElementById(`checkbox-${id}`) as HTMLInputElement;
      return el ? el.checked : false;
    },
    // Shortcut for getting blank field value.
    blankValue: () => {
      let el = document.getElementById('message-dialog-blank') as HTMLInputElement;
      return el ? el.value : '';
    }
  };
}
export const messageDialog = createMessage();
