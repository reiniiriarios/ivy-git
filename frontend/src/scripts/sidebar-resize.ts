let sidebar: HTMLElement;
let main: HTMLElement;

const minSidebar = 300;
const maxSidebar = 800;
const defaultSidebarWidth = document.documentElement.style.getPropertyValue('--sidebar-width');

export const setMainBlock = (el: HTMLElement) => {
  main = el;
}

export const resizableSidebar = (el: HTMLElement) => {
  sidebar = el;

  // Starting x-coord on click.
  let x = 0;
  // Starting width of sidebar on click.
  let sw = 0;
  // Starting width of main on click.
  let mw = 0;

  const resizer = document.createElement('div');
  resizer.classList.add('resizer-v');
  resizer.style.height = sidebar.offsetHeight + 'px';

  const mouseDownHandler = (e: MouseEvent) => {
    // Prevent accidentally scrolling the commits list.
    main.classList.add('dragging');
    sidebar.classList.add('dragging');

    // Current x-coord of mouse.
    x = e.pageX;
    sw = parseInt(window.getComputedStyle(sidebar).width);
    mw = parseInt(window.getComputedStyle(main).width);

    document.addEventListener('mousemove', mouseMoveHandler);
    document.addEventListener('mouseup', mouseUpHandler);
  }

  const mouseMoveHandler = (e: MouseEvent) => {
    let move = e.pageX - x;
    let swn = Math.min(Math.max((sw + move), minSidebar), maxSidebar);
    document.documentElement.style.setProperty('--sidebar-width', swn + 'px');
  }

  const mouseUpHandler = () => {
    main.classList.remove('dragging');
    sidebar.classList.remove('dragging');

    document.removeEventListener('mousemove', mouseMoveHandler);
    document.removeEventListener('mouseup', mouseUpHandler);
  }

  resizer.addEventListener('mousedown', mouseDownHandler);

  sidebar.appendChild(resizer);
}
