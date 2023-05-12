let commitsContainer: HTMLElement;
let commitDetailsContainer: HTMLElement;

let minCommits = 200;
let minDetails = 100;

export const setCommitsContainer = (el: HTMLElement) => {
  commitsContainer = el;
}

export const setDetailsResizable = (el: HTMLElement) => {
  commitDetailsContainer = el;

  // Starting y-coord on click.
  let y = 0;
  // Starting height of details container on click.
  let dh = 0;
  // Starting height of commits container on click.
  let ch = 0;
  let max = 0;

  const resizer = document.createElement('div');
  resizer.classList.add('resizer-h');
  resizer.style.width = commitDetailsContainer.offsetWidth + 'px';

  const mouseDownHandler = (e: MouseEvent) => {
    // Prevent accidentally scrolling the commits list.
    commitsContainer.style.pointerEvents = 'none';

    // Current x-coord of mouse.
    y = e.pageY;
    dh = parseInt(window.getComputedStyle(commitDetailsContainer).height);
    ch = parseInt(window.getComputedStyle(commitsContainer).height);

    // Lock commits window at 200px tall.
    max = dh + ch - minCommits;

    document.addEventListener('mousemove', mouseMoveHandler);
    document.addEventListener('mouseup', mouseUpHandler);
  }

  const mouseMoveHandler = (e: MouseEvent) => {
    let move = y - e.pageY;
    let dhn = Math.min(Math.max((dh + move), minDetails), max);
    let chn = ch + dh - dhn;
    commitDetailsContainer.style.height = dhn + 'px';
    commitsContainer.style.height = chn + 'px';
  }

  const mouseUpHandler = () => {
    commitsContainer.style.pointerEvents = 'all';

    document.removeEventListener('mousemove', mouseMoveHandler);
    document.removeEventListener('mouseup', mouseUpHandler);
  }

  resizer.addEventListener('mousedown', mouseDownHandler);
  window.addEventListener('resize', makeDetailsSizingPercentage);

  commitDetailsContainer.appendChild(resizer);
}

// Resize window handler. Update as a percentage with min height.
const makeDetailsSizingPercentage = async (e: UIEvent) => {
  let ch = parseInt(window.getComputedStyle(commitsContainer).height);
  let dh = parseInt(window.getComputedStyle(commitDetailsContainer).height);
  if (ch < minCommits) ch = minCommits;
  if (dh < minDetails) dh = minDetails;
  let chp = ch / (ch + dh);
  let dhp = dh / (ch + dh);
  commitsContainer.style.height = (chp * 100).toFixed(4) + "%";
  commitDetailsContainer.style.height = (dhp * 100).toFixed(4) + "%";
}
