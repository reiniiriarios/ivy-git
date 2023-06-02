import octicons from "@primer/octicons";

function checkIcon(el: HTMLElement, reenable: boolean = false) {
  let icon = el.getElementsByClassName('icon')[0];
  let normalSvg = icon.innerHTML;
  icon.innerHTML = octicons['check'].toSVG({ "width": 14 });
  icon.classList.add('success');
  setTimeout(() => {
    icon.classList.remove('success');
    icon.innerHTML = normalSvg;
    if (reenable) {
      el.removeAttribute('disabled');
    }
  }, 1000);
}

export default checkIcon;
