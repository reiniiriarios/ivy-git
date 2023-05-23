import type { EnvironmentInfo } from "wailsjs/runtime/runtime";

let platform: string;

export const envInit = async () => {
  (window as any).runtime.Environment().then((env: EnvironmentInfo) => {
    platform = env.platform;
    //...
  });
}

export function isDarwin(): boolean {
  return platform === 'darwin';
}

export function getPlatform(): string {
  return platform;
}
