let tableElement: HTMLTableElement;
let tableCols: Col[] = [];

interface Col {
  e: HTMLElement;
  minWidth: number;
}

export const setCommitsTable = (e: HTMLTableElement) => {
  tableElement = e;
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

  if (col.dataset.name === 'branch') {
    min = parseInt(window.getComputedStyle(col).width);
  } else if (col.dataset.name === 'tree') {
    min = 50;
  } else {
    min = 75;
  }
  col.style.minWidth = min + 'px';

  tableCols[parseInt(col.dataset.order)] = {
    e: col,
    minWidth: min,
  };

  if (typeof col.dataset.nograb === "undefined") {
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
      col2.e.style.width = (w2 + (w - newWidth)) + 'px';
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

const updateSavedTableSizing = async () => {
  let tableWidth = parseInt(window.getComputedStyle(tableElement).width);

  // Calculate all columns, convert all to percentage to keep when resizing.
  for (let i = 0; i < tableCols.length; i++) {
    let colWidth = parseInt(window.getComputedStyle(tableCols[i].e).width);
    let newWidth = colWidth / tableWidth;
    if (colWidth < tableCols[i].minWidth) {
      newWidth = tableCols[i].minWidth / tableWidth;
    }
    tableCols[i].e.style.width = (newWidth * 100).toFixed(4) + '%';
  }

  console.log(tableCols);
}

const resizeWindow = async (e: UIEvent) => {
  console.log(e);
  updateSavedTableSizing();
}
