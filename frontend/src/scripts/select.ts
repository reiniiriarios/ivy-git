// Attach to all select elements with <select use:select>.
export function select(el: HTMLSelectElement) {
  // Build
  let select = document.createElement('div');
  select.classList.add('select');
  let button = document.createElement('button');
  let ul = document.createElement('ul');

  let options = el.getElementsByTagName('option');
  let lis: HTMLLIElement[] = [];
  for (let i = 0; i < options.length; i++) {
    if (el.dataset.required === 'true' && !options[i].value) {
      if (options[i].selected) {
        button.innerText = options[i].innerText;
      }
      continue;
    }
    let li = document.createElement('li');
    li.innerText = options[i].textContent;
    li.dataset.value = options[i].value;
    if (options[i].selected) {
      button.innerText = li.innerText;
      li.classList.add('selected');
    }
    lis.push(li);
    ul.appendChild(li);
  }

  select.appendChild(button);
  select.appendChild(ul);

  // Events
  function click(e: MouseEvent | KeyboardEvent) {
    e.preventDefault();
    select.classList.toggle('open');

    if (e.target instanceof HTMLLIElement) {
      el.value = e.target.dataset.value;
      button.innerText = e.target.innerText;
      console.log(lis);
      for (let i = 0; i < lis.length; i++) {
        lis[i].classList.remove('selected');
      }
      e.target.classList.add('selected');
      el.dispatchEvent(new CustomEvent('change'));
    }
  }

  select.addEventListener('click', click);
  select.addEventListener('keypress', click);

  // Replace
  el.parentNode.insertBefore(select, el);
  el.style.display = 'none';

  // Hide if clicking elsewhere
  document.addEventListener('click', function(e) {
    if (!select.contains(e.target as any)) {
      select.classList.remove('open');
    }
  });
}
