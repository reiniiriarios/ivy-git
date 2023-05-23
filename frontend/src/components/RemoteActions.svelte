<script lang="ts">
  import octicons from "@primer/octicons";
  import { parseResponse } from "scripts/parse-response";
  import { currentRemote, remoteData } from "stores/remotes";
  import { onMount } from 'svelte';
  import { FetchRemote, PullRemote, PushRemote } from "wailsjs/go/main/App";

  function pull(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    PullRemote(remote).then((r) => {
      parseResponse(r, () => {
        checkIcon(el);
        remoteData.refresh();
      }, () => el.removeAttribute('disabled'));
    });
  }

  function push(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    PushRemote(remote).then((r) => {
      parseResponse(r, () => {
        checkIcon(el);
        remoteData.refresh();
      }, () => el.removeAttribute('disabled'));
    });
  }

  function fetch(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    FetchRemote(remote).then((r) => {
      parseResponse(r, () => {
        checkIcon(el, true);
        remoteData.refresh();
      }, () => el.removeAttribute('disabled'));
    });
  }

  function checkIcon(el: HTMLElement, reenable: boolean = false) {
    let icon = el.getElementsByClassName('icon')[0];
    let normalSvg = icon.innerHTML;
    icon.innerHTML = octicons['check'].toSVG({ "width": 14 });
    icon.classList.add('success');
    setTimeout(() => {
      icon.classList.remove('success');
      icon.innerHTML = normalSvg;
      if (reenable) {
        el.removeAttribute('disabled');
      }
    }, 1000);
  }
</script>

<!-- todo: single sync action -->
<div class="sidebar-remote-actions">
  <button class="btn" on:click={e => fetch(e, $currentRemote.Name)}>
    <span class="icon">{@html octicons['arrow-down-left'].toSVG({ "width": 14 })}</span>
    Fetch
  </button>
  <button class="btn" on:click={e => push(e, $currentRemote.Name)} disabled={$currentRemote.Ahead == 0 || !$currentRemote.Push}>
    <span class="icon">{@html octicons['arrow-up'].toSVG({ "width": 14 })}</span>
    Push
  </button>
  <button class="btn" on:click={e => pull(e, $currentRemote.Name)} disabled={$currentRemote.Behind == 0 || !$currentRemote.Fetch}>
    <span class="icon">{@html octicons['arrow-down'].toSVG({ "width": 14 })}</span>
    Pull
  </button>
</div>
