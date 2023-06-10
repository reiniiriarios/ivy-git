<script lang="ts">
  export let display: string = "";
  export let id: string = "";
  export let classes: string = "";
  export let value: any;
  export let valid: boolean = true;
  export let use: (el: HTMLInputElement) => void = () => {};
  export let validate: (value: any) => boolean = null;

  $: if (validate) valid = validate(value);
</script>

<label class="text-input {classes}">
  {#if display}
    <span>{display}</span>
  {/if}
  <input type="text"
    use:use
    {id}
    bind:value={value}
    class:invalid={value && !valid}
    on:input={() => {
      if (validate) valid = validate(value)
    }}
  >
</label>
