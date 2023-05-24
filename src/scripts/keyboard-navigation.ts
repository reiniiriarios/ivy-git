import { isDarwin } from "scripts/env";

export function keyboardNavListener() {
  window.addEventListener('keydown', (e: KeyboardEvent) => {
    if (!['input', 'select', 'textarea'].includes(document.activeElement?.tagName.toLowerCase())) {
      if(['ArrowDown', 'ArrowUp', 'ArrowLeft', 'ArrowRight'].includes(e.key)) {
        e.preventDefault();
        if (document.activeElement) {
          let es = document.querySelectorAll(
            'a:not([disabled]), button:not([disabled]), input[type=text]:not([disabled]), [tabindex]:not([disabled]):not([tabindex="-1"])'
          ) as NodeListOf<HTMLElement>;
          let elems = [...es].filter(e => e.offsetWidth > 0 || e.offsetHeight > 0 || e === document.activeElement);
          let i = elems.indexOf(document.activeElement as HTMLElement);
          if (i > -1) {
            if (['ArrowDown', 'ArrowRight'].includes(e.key)) {
              if (elems[i + 1]) {
                elems[i + 1].focus();
              }
            } else {
              if (elems[i - 1]) {
                elems[i - 1].focus();
              }
            }
          }
        }
      }
    }
  });
}

// Something is breaking Cmd/Ctrl+A to select all text within input elements.
// This resolves that issue, but is hacky / should be investigated further.
export function addInputListener() {
  document.addEventListener('keydown', (e: KeyboardEvent) => {
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
      let cmd = (isDarwin() && e.metaKey) || (!isDarwin() && e.ctrlKey);
      if (cmd && (e.key === 'A' || e.key === 'a')) {
        e.target.select();
      }
    }
  })
}
