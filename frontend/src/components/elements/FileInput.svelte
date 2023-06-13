<script lang="ts">
  import { parseResponse } from "scripts/parse-response";
  import { SelectDirectory } from "wailsjs/go/main/App";

  export let display: string = "";
  export let id: string = "";
  export let classes: string = "";
  export let value: any;
  export let valid: boolean = true;
  export let use: (el: HTMLInputElement) => void = () => {};
  export let directory: boolean = false;

  function choose() {
    if (directory) {
      SelectDirectory().then(result => {
        parseResponse(result, () => {
          if (result.Response !== 'none') {
            value = result.Data;
          }
        });
      });
    }
    else {
      // possibly todo
      console.log('file selection not implemented');
    }
  }
</script>

<label class="file-input {classes}">
  {#if display}
    <span class="file-input__display">{display}</span>
  {/if}
  <input type="text"
    class="file-input__input"
    use:use
    {id}
    bind:value={value}
    class:invalid={value && !valid}
    on:click={choose}
    on:change
    readonly
  >
  <button class="btn file-input__btn" on:click={choose}>Choose</button>
</label>
