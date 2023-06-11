import { parseResponse } from 'scripts/parse-response';
import { writable } from 'svelte/store';
import { GetGitConfigGlobal, GetGitConfigLocal } from 'wailsjs/go/main/App';

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
    set: async (value: any) => {
      set(value);
    },
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
  };
}
export const gitConfig = createGitConfig();
