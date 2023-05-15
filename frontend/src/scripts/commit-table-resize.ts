let commitsElement: HTMLElement;
let tableElement: HTMLTableElement;
let tableCols: Col[] = [];

interface Col {
  e: HTMLElement;
  minWidth: number;
  flex: boolean;
}

export const setCommitsTable = (e: HTMLTableElement) => {
  tableElement = e;
  // This allows auto table layout to set the widths,
  // those widths are saved as distinct pixel sizes,
  // then the table is converted to a fixed layout.
  updateSavedTableSizing().then(() => {
    tableElement.style.tableLayout = 'fixed';
  });
  window.addEventListener('resize', resizeWindow);
}

export const createResizableColumn = (col: HTMLElement) => {
  // Starting x-coord on click.
  let x = 0;
  // Starting width of column on click.
  let w = 0;
  // Column to the right.
  let col2: Col;
  // Starting width of column to the right on click.
  let w2 = 0;
  // Minimum width this column can be.
  let min = 0;
  // Maximum width this column can be, based on column to the right.
  let max = 0;
  // Columns marked to flex are not given specific widths and cannot resize on their own.
  let flex = typeof col.dataset.resizeflex !== "undefined";

  // Get min width for all columns.
  if (col.dataset.name === 'branch') {
    min = parseInt(window.getComputedStyle(col).width);
  } else {
    min = 75;
  }
  // The tree column is already set separately by the graph.
  if (col.dataset.name !== 'tree') {
    col.style.minWidth = min + 'px';
  }

  tableCols[parseInt(col.dataset.order)] = {
    e: col,
    minWidth: min,
    flex: flex,
  };

  if (!flex) {
    const resizer = document.createElement('div');
    resizer.classList.add('resizer');
    resizer.style.height = col.offsetHeight + 'px';

    const mouseDownHandler = (e: MouseEvent) => {
      // Current x-coord of mouse.
      x = e.pageX;
      w = parseInt(window.getComputedStyle(col).width);
      // Column to the right.
      col2 = tableCols[parseInt(col.dataset.order) + 1];
      w2 = parseInt(window.getComputedStyle(col2.e).width);

      // This max width prevents this column resizing the column two columns to the right.
      max = w + w2 - col2.minWidth;

      document.addEventListener('mousemove', mouseMoveHandler);
      document.addEventListener('mouseup', mouseUpHandler);
    }

    const mouseMoveHandler = (e: MouseEvent) => {
      let newWidth = Math.min(Math.max((w + e.pageX - x), min), max);
      // Resize both this column and the column to the right in equal amounts.
      col.style.width = newWidth + 'px';
      if (!col2.flex) {
        col2.e.style.width = (w2 + (w - newWidth)) + 'px';
      }
    }

    const mouseUpHandler = () => {
      updateSavedTableSizing();

      document.removeEventListener('mousemove', mouseMoveHandler);
      document.removeEventListener('mouseup', mouseUpHandler);
    }

    resizer.addEventListener('mousedown', mouseDownHandler);

    col.appendChild(resizer);
  }
}

// Update all columns, setting their current width to their current computed width.
export const updateSavedTableSizing = async () => {
  for (let i = 0; i < tableCols.length; i++) {
    if (!tableCols[i].flex) {
      tableCols[i].e.style.width = window.getComputedStyle(tableCols[i].e).width;
    }
  }
}

// Update all columns, setting their current width to their current computed width.
export const removeTableSizing = async () => {
  if (!tableElement) tableElement = document.getElementById('commits__table') as HTMLTableElement;
  if (tableElement) tableElement.style.tableLayout = 'auto';
  if (!commitsElement) commitsElement = document.getElementById('commits__scroll') as HTMLElement;
  if (commitsElement) commitsElement.scrollTop = 0;
  for (let i = 0; i < tableCols.length; i++) {
    if (!tableCols[i].flex) {
      tableCols[i].e.style.width = '';
    }
  }
}

// Resize window handler.
// Update all columns, setting their current width to their current computed width, as a percentage.
const resizeWindow = async (e: UIEvent) => {
  let tableWidth = parseInt(window.getComputedStyle(tableElement).width);

  for (let i = 0; i < tableCols.length; i++) {
    if (!tableCols[i].flex) {
      let colWidth = parseInt(window.getComputedStyle(tableCols[i].e).width);
      // If column larger than min width, limit to min width.
      let newWidth =
        colWidth < tableCols[i].minWidth
          ? tableCols[i].minWidth / tableWidth
          : colWidth / tableWidth;
      tableCols[i].e.style.width = (newWidth * 100).toFixed(4) + "%";
    }
  }
}
