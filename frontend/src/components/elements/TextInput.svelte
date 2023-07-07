<script lang="ts">
  import Info from "./Info.svelte";

  // Text do display for input label:
  export let display: string = "";
  // Input `id`.
  export let id: string = "";
  // Input `class`.
  export let classes: string = "";
  // Input value.
  export let value: any;
  // Input is valid, based on validate().
  export let valid: boolean = true;

  // Passes `use` up.
  export let use: (el: HTMLInputElement) => void = () => {};
  // Whether the value is valid. Styles input.
  export let validate: (value: any) => boolean = null;
  // If the input value changes, this method will auto-edit it.
  // Use to auto-remove some values, or edit them.
  export let autoEditBlank: (value: string) => string = null;
  // Runs on:change.
  export let update: (var1: any, var2?: any) => void = null;

  // `Info` component:
  export let tip: string = "";
  export let tipPosition: string = "right";

  $: if (autoEditBlank) value = autoEditBlank(value);
  $: if (validate) valid = validate(value);
</script>

<label class="text-input {classes}">
  {#if display}
    <span class="text-input__display">
      {display}
      {#if tip}
        <Info position={tipPosition}>{@html tip}</Info>
      {/if}
    </span>
  {/if}
  <input type="text"
    use:use
    {id}
    bind:value={value}
    class:invalid={value && !valid}
    on:input={() => {
      if (autoEditBlank) value = autoEditBlank(value);
      if (validate) valid = validate(value);
    }}
    on:change={() => {
      if (update) update(value);
    }}
  >
</label>
