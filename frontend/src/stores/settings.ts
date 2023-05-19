import { writable } from 'svelte/store';
import { GetSettings } from 'wailsjs/go/main/App'

interface Settings {
	Version: string;
	DisplayCommitSignatureInList: boolean;
}

// These stores reflec the current ui state and can be used
// across the app to change the ui state from components
// unrelated in hierarchy, but related in content.
function createSettings() {
  const { subscribe, set } = writable({} as Settings);
  
  return {
    subscribe,
    refresh: async () => {
      GetSettings().then(result => {
        set(result);
      });
    },
  };
}
export const settings = createSettings();
