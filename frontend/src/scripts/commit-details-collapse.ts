export function toggleDir(e: MouseEvent & { currentTarget: HTMLElement } | KeyboardEvent & { currentTarget: HTMLElement }) {
  let el = e.currentTarget.parentNode.parentNode as HTMLElement;
  el.classList.contains('filestatdir__dir--closed') ? expandDir(el) : collapseDir(el);
}

export function resetDirs() {
  let els = document.getElementsByClassName('filestatdir__dir--closed') as HTMLCollectionOf<HTMLElement>;
  for (let i = 0; i < els.length; i++) {
    els[i].classList.remove('filestatdir__dir--closed');
  }
}

function collapseDir(el: HTMLElement) {
  el.classList.add('filestatdir__dir--closed');
}

function expandDir(el: HTMLElement) {
  el.classList.remove('filestatdir__dir--closed');
}
