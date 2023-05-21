import { writable } from 'svelte/store';
import { GetRemotes } from 'wailsjs/go/main/App';

export interface Remote {
	Name: string;
	Url: string;
	Fetch: boolean;
	Push: boolean;
	Site: string;
	Repo: string;
	User: string;
	RepoName: string;
	Ahead: number;
	Behind: number;
	LastUpdate: number;
}

function createRemotes() {
  const { subscribe, set } = writable([] as Remote[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetRemotes().then(result => {
        console.log(result);
        set(result.Data as Remote[]);
      });
    },
  };
}
export const remotes = createRemotes();
