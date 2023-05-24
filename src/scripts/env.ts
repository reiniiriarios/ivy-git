type EnvironmentInfo = any;

let environment: EnvironmentInfo;

export const envInit = async (): Promise<EnvironmentInfo> => {
  // await (window as any).runtime.Environment().then((env: EnvironmentInfo) => environment = env);
  return {
    platform: 'darwin',
  };
}

export function isDarwin(): boolean {
  return environment.platform === 'darwin';
}

export function getPlatform(): string {
  return environment.platform;
}
