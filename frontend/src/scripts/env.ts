import type { EnvironmentInfo } from "wailsjs/runtime/runtime";

let environment: EnvironmentInfo;
let platform: string;

export const envInit = async (): Promise<EnvironmentInfo> => {
  await (window as any).runtime.Environment().then((env: EnvironmentInfo) => environment = env);
  return environment;
}

export function isDarwin(): boolean {
  return platform === 'darwin';
}

export function getPlatform(): string {
  return platform;
}
