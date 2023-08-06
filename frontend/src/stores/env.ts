import { derived, get, writable } from "svelte/store";
import type { EnvironmentInfo } from "wailsjs/runtime/runtime";

function createEnv() {
  const { subscribe, set } = writable({} as EnvironmentInfo);
  
  return {
    subscribe,
    fetch: async () => {
      await (window as any).runtime.Environment().then((env: EnvironmentInfo) => set(env));
      return get(environment);
    }
  }
}
export const environment = createEnv();
export const isDarwin = derived(environment, $environment => $environment?.platform === 'darwin');
export const finderWord = derived(environment, $environment => {
  switch ($environment?.platform) {
    case 'darwin':
      return 'Finder';
    case 'windows':
      return 'Explorer';
    default:
      return 'File Browser';
  }
});
