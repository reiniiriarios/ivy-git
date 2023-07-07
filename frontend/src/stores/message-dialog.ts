import { writable } from 'svelte/store';

type MessageCallback = (e?: MouseEvent & { currentTarget: HTMLElement } | KeyboardEvent & { currentTarget: HTMLElement }) => any;

interface Message {
  // Basic values:
  heading?: string;
  message?: string;
  confirm?: string;
  okay?: string;

  // Yes/No callbacks:
  callbackConfirm?: MessageCallback;
  callback?: MessageCallback;

  // Forms:
  options?: {
    text: string;
    icon?: string;
    callback: MessageCallback;
  }[];
  checkboxes?: {
    id: string;
    label: string;
    checked?: boolean;
  }[];
  blank?: string;
  validateBlank?: (input: string) => boolean,
  autoEditBlank?: (input: string) => string,

  // Specific message dialogs:
  addTag?: boolean;
  newRepo?: boolean;
  cloneRepo?: boolean;
  addRemote?: boolean;
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
        blank: message.blank ?? '',
        validateBlank: message.validateBlank ?? null,
        autoEditBlank: message.autoEditBlank ?? null,
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    addTag: async(message: Message) => {
      set({
        heading: message.heading ?? 'Add Tag',
        message: message.message ?? 'Enter data:',
        confirm: message.confirm ?? 'Add',
        addTag: true,
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    addRepo: async(message: Message) => {
      set({
        heading: message.heading ?? 'Create New Repo',
        message: message.message ?? 'Enter a name and select a directory to create a new repository in.',
        confirm: message.confirm ?? 'Create',
        newRepo: true,
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    cloneRepo: async(message: Message) => {
      set({
        heading: message.heading ?? 'Clone Repo',
        message: message.message ?? 'Enter a URL and select a directory to clone the repository into.',
        confirm: message.confirm ?? 'Clone',
        cloneRepo: true,
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    addRemote: async(message: Message) => {
      set({
        heading: message.heading ?? 'Add Remote',
        message: message.message ?? 'Enter a name for the remote and URL(s).',
        confirm: message.confirm ?? 'Add',
        addRemote: true,
        callbackConfirm: message.callbackConfirm ?? (() => {}),
        okay: message.okay ?? 'Cancel',
        callback: message.callback ?? (() => {}),
      });
    },
    yes: async() => {
      update(message => {
        if (message.callbackConfirm) {
          message.callbackConfirm();
        }
        return {};
      });
    },
    okay: async() => {
      update(message => {
        if (message.callback) {
          message.callback();
        }
        return {};
      });
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
    },
    // Shortcut for getting add tag data.
    addTagData: () => {
      return {
        name: (document.getElementById('message-dialog-tag-name') as HTMLInputElement).value,
        message: (document.getElementById('message-dialog-tag-message') as HTMLInputElement).value,
        push: (document.getElementById('message-dialog-tag-push') as HTMLInputElement).checked,
      }
    },
    // Shortcut for getting add repo data.
    addRepoData: () => {
      return {
        name: (document.getElementById('message-dialog-repo-name') as HTMLInputElement).value,
        location: (document.getElementById('message-dialog-repo-location') as HTMLInputElement).value,
      }
    },
    // Shortcut for getting clone repo data.
    cloneRepoData: () => {
      return {
        url: (document.getElementById('message-dialog-repo-url') as HTMLInputElement).value,
        location: (document.getElementById('message-dialog-repo-location') as HTMLInputElement).value,
      }
    },
    // Shortcut for getting add remote data.
    addRemoteData: () => {
      return {
        name: (document.getElementById('message-dialog-remote-name') as HTMLInputElement).value,
        fetch: (document.getElementById('message-dialog-remote-fetch') as HTMLInputElement).value,
        push: (document.getElementById('message-dialog-remote-push') as HTMLInputElement).value,
      }
    },
  };
}
export const messageDialog = createMessage();
