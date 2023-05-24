import { derived, writable } from 'svelte/store';
import { GetRemotes } from 'src/_tmp';

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
  const { subscribe, set } = writable({
    Remotes: [] as Remote[],
    CurrentRemote: "",
  });
  
  return {
    subscribe,
    refresh: async () => {
      GetRemotes().then(result => {
        console.log(result);
        set(result.Data);
      });
    },
  };
}
export const remoteData = createRemotes();
export const remotes = derived(remoteData, $remoteData => $remoteData.Remotes);
export const currentRemote = derived(remoteData, $remoteData => $remoteData.Remotes.find(r => r.Name === $remoteData.CurrentRemote));
