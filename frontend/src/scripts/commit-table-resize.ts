export const createResizableColumn = (col: HTMLElement) => {
  let x = 0;
  let w = 0;

  console.log(col);

  const resizer = document.createElement('div');
  resizer.classList.add('resizer');
  resizer.style.height = col.offsetHeight + 'px';

  const mouseDownHandler = (e: MouseEvent) => {
    x = e.pageX;
    w = parseInt(window.getComputedStyle(col).width);

    document.addEventListener('mousemove', mouseMoveHandler);
    document.addEventListener('mouseup', mouseUpHandler);
  }

  const mouseMoveHandler = (e: MouseEvent) => {
    col.style.width = (w + e.pageX - x) + 'px';
  }

  const mouseUpHandler = () => {
    document.removeEventListener('mousemove', mouseMoveHandler);
    document.removeEventListener('mouseup', mouseUpHandler);
  }

  resizer.addEventListener('mousedown', mouseDownHandler);

  col.appendChild(resizer);
}
