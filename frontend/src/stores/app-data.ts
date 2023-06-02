import { writable } from 'svelte/store';
import { GetAppData } from 'wailsjs/go/main/App';
import { EventsOn } from 'wailsjs/runtime/runtime';

interface AppData {
	RecentRepoDir: string;
}

function createAppData() {
  const { subscribe, set } = writable({} as AppData);

  EventsOn('appdata', (newAppData: AppData) => {
    appData.set(newAppData);
  });

  return {
    subscribe,
    set,
    fetch: async () => {
      GetAppData().then(result => set(result));
    },
  };
}
export const appData = createAppData();
