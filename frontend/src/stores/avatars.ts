import { get, writable } from "svelte/store";
import { GetAvatarUrl } from "wailsjs/go/main/App";

type Avatars = { [email: string]: string };

// This store acts as a cache for avatar urls and is repopulated every app restart.
function createAvatars() {
  const { subscribe, update } = writable({} as Avatars);

  return {
    subscribe,
    fetch: async (email: string): Promise<string> => {
      if (!get(avatars)[email]) {
        await GetAvatarUrl(email).then(url => {
          update(avatars => {
            avatars[email] = url;
            return avatars;
          });
        });
      }
      return get(avatars)[email];
    },
  };
}
export const avatars = createAvatars();
