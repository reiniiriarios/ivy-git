<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

  export let values: string[] | number[] = [];
  export let options: { [value: string|number]: string } = {};
  export let selected: string|number;

  let select: HTMLDivElement;

  let currentValue: string;
  let currentDisplay: string = "—";
  $: currentValue = selected.toString();

  let opts: { value: string|number, display: string }[] = [];
  $: {
    if (!options.length && values.length) {
      opts = values.map((v: string|number) => {
        if (v === currentValue) currentDisplay = v.toString();
        return {
          value: v,
          display: v.toString(),
        };
      });
    } else {
      opts = [];
      Object.keys(options).forEach(v => {
        if (v === currentValue) currentDisplay = options[v];
        opts.push({
          value: v,
          display: options[v],
        });
      });
    }
  }

  function toggleOpen() {
    select.classList.toggle('open');
  }

  function selectOption(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    currentValue = e.currentTarget.dataset.value;
    currentDisplay = e.currentTarget.innerText;
		dispatch('change', {
			value: currentValue,
      display: currentDisplay,
		});
  }

  document.addEventListener('click', function(e) {
    if (select && !select.contains(e.target as any)) {
      select.classList.remove('open');
    }
  });
</script>

<div class="select-container">
  <div bind:this={select} class="select" on:click={toggleOpen} on:keypress={toggleOpen}>
    <button>{currentDisplay}</button>
    <ul>
      {#each opts as {value, display}}
        <li class:selected={currentValue === value} data-value="{value}" on:click={selectOption} on:keypress={selectOption}>{display}</li>
      {/each}
    </ul>
  </div>
</div>
