export function checkRef(ref: string): boolean {
  // Must be at least one character.
  if (!ref) return false;
  // Cannot be "@"
  if (ref === "@") return false;
  // Cannot contain the following sequences.
  if (ref.includes("/.") || ref.includes("//") || ref.includes("@{") || ref.includes("..") || ref.includes(".lock/")) return false;
  // Cannot start with.
  if (ref[0] === '/' || ref[0] === '.') return false;
  // Cannot end with.
  if (ref[ref.length - 1] === '/' || ref[ref.length - 1] === '.' || ref.substring(ref.length - 5) === '.lock') return false;
  // Cannot contain the following characters, inc ascii control characters or DEL.
  if (ref.match(/[\x00-\x20\x1f ~\^:\?\*\[\\]/)) return false;

  return true;
}
