export function tabUpDown(e: KeyboardEvent) {
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
