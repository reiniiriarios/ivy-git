<script lang="ts">
  import Info from "./Info.svelte";

  export let display: string = "";
  export let id: string = "";
  export let classes: string = "";
  export let value: any;
  export let valid: boolean = true;
  export let use: (el: HTMLInputElement) => void = () => {};
  export let validate: (value: any) => boolean = null;
  export let update: (var1: any, var2?: any) => void = null;
  export let tip: string = "";
  export let tipPosition: string = "right";

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
      if (validate) valid = validate(value);
    }}
    on:change={() => {
      if (update) update(value);
    }}
  >
</label>
