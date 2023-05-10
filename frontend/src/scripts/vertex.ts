export const vertexOver = (e: MouseEvent & { target: SVGGElement }) => {
  let v = e.target.dataset.id;
  console.log(e);
  getCommitRow(v).classList.add('hover');
}

export const vertexOut = (e: MouseEvent & { target: SVGGElement }) => {
  let v = e.target.dataset.id;
  console.log(e);
  getCommitRow(v).classList.remove('hover');
}

const getCommitRow = (id: string): HTMLElement => {
  return document.querySelector(`.commit[data-id="${id}"]`);
}
