// Validate GPG Key Format
export function checkGpgKeyFormat(key: string): boolean {
  // empty is okay, that's valid
  if (!key) return true;
  key = key.replaceAll(' ','');
  // 8, 16, or 40 digits
  return /^(?:(?:[A-F0-9]{24})?[A-F0-9]{8})?[A-F0-9]{8}$/i.test(key);
}

// Clean GPG Key to last 8 digits.
export function cleanGpgKey(key: string): string {
  key = key.replaceAll(' ','');
  return key.substring(key.length - 8);
}
