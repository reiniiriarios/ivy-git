import { derived, writable } from 'svelte/store';
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
  const { subscribe, set } = writable({
    Remotes: [] as Remote[],
    CurrentRemote: "",
  });
  
  return {
    subscribe,
    refresh: async () => {
      GetRemotes().then(result => {
        // Ignore errors here, there may be no remotes.
        set(result.Data);
      });
    },
  };
}
export const remoteData = createRemotes();
export const remotes = derived(remoteData, $remoteData => {
  console.log($remoteData);
  return $remoteData?.Remotes?.length ? $remoteData.Remotes : [];
});
export const currentRemote = derived(remoteData, $remoteData => {
  if (!$remoteData?.Remotes?.length) return {} as Remote;
  return $remoteData.Remotes.find(r => r.Name === $remoteData.CurrentRemote) ?? {} as Remote;
});
