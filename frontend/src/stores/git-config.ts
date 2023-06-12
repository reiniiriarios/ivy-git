import { parseResponse } from 'scripts/parse-response';
import { get, writable } from 'svelte/store';
import { GetGitConfigGlobal, GetGitConfigLocal, UpdateGitConfigSignCommits, UpdateGitConfigUserEmail, UpdateGitConfigUserName, UpdateGitConfigUserSigningKey } from 'wailsjs/go/main/App';

interface GitConfigAll {
  local: GitConfig,
  global: GitConfig,
  // system
}

interface GitConfig {
	UserName: string;
	UserEmail: string;
	UserSigningKey: string;
	CommitGpgSign: boolean;
}

function createGitConfig() {
  const { subscribe, set, update } = writable({
    local: {} as GitConfig,
    global: {} as GitConfig,
  } as GitConfigAll);
  
  return {
    subscribe,
    set,
    fetch: async () => {
      GetGitConfigLocal().then(result => {
        parseResponse(result, () => {
          update(cfg => {
            cfg.local = result.Data;
            return cfg;
          });
        });
      });
      GetGitConfigGlobal().then(result => {
        parseResponse(result, () => {
          update(cfg => {
            cfg.global = result.Data;
            return cfg;
          });
        });
      });
    },
    updateSetting: async (list: string, setting: string, value: string) => {
      update(cfg => {
        cfg[list][setting] = value;
        return cfg;
      });
    },
    setUserName: (list: string, value: string) => {
      UpdateGitConfigUserName(list, value).then(r => parseResponse(r));
    },
    setUserEmail: (list: string, value: string) => {
      UpdateGitConfigUserEmail(list, value).then(r => parseResponse(r));
    },
    setUserSigningKey: (list: string, value: string) => {
      UpdateGitConfigUserSigningKey(list, value).then(r => parseResponse(r));
    },
    setSignCommits: (list: string, value: boolean) => {
      UpdateGitConfigSignCommits(list, value).then(r => parseResponse(r));
    },
  };
}
export const gitConfig = createGitConfig();
