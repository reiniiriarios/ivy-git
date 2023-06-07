const conventionalLabels = [
  'feat',
  'fix',
  'refactor',
  'chore',
  'docs',
  'style',
  'test',
  'revert',
  'drop',
  'wip',
];

export function highlightConventionalCommits(subject: string): string {
  let r = /^(([a-z]+)(?:\(.+\))?:)(.*)$/i
  let matches = subject.match(r);
  if (matches?.length === 4 && conventionalLabels.includes(matches[2].toLowerCase())) {
    let label = matches[2].toLowerCase();
    // $0 all
    // $1 label(scope):
    // $2 label
    // $3 everything else
    return subject.replace(r, `<span class='conv-commit-label conv-commit-label--${label}'>$1</span>$3`);
  }
  return subject;
}
