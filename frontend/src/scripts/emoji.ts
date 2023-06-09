import emoji from "./emoji.json";

export function emojify(string: string): string {
  let r = /(:[a-z0-9-]+:)/ig
  let matches = string.match(r);
  if (matches?.length === 4 && typeof emoji[matches[2].toLowerCase()] !== 'undefined') {
    return string.replaceAll(r, emoji[matches[2].toLowerCase()]);
  }
  return string;
}
