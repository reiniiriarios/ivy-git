<script lang="ts">
  import octicons from "@primer/octicons";
  import fetchRemote from "actions/remote/fetch";
  import pullRemote from "actions/remote/pull";
  import pushRemote from "actions/remote/push";
  import { currentRemote } from "stores/remotes";
  import { branchSelect, repoSelect } from "stores/ui";

  function pull(e: MouseEvent | KeyboardEvent, remote: string) {
    pullRemote(remote, e.target as HTMLElement);
  }

  function push(e: MouseEvent | KeyboardEvent, remote: string) {
    pushRemote(remote, e.target as HTMLElement);
  }

  function fetch(e: MouseEvent | KeyboardEvent, remote: string) {
    fetchRemote(remote, e.target as HTMLElement);
  }
</script>

<!-- todo: single sync action -->
<div class="sidebar-remote-actions" style:display={$repoSelect || $branchSelect ? 'none' : 'flex'}>
  <button class="btn btn--icon" on:click={e => fetch(e, $currentRemote.Name)}>
    <span class="icon">{@html octicons['arrow-down-left'].toSVG({ "width": 14 })}</span>
    Fetch
  </button>
  <button class="btn btn--icon" on:click={e => push(e, $currentRemote.Name)} disabled={$currentRemote.Ahead == 0 || !$currentRemote.Push}>
    <span class="icon">{@html octicons['arrow-up'].toSVG({ "width": 14 })}</span>
    Push
  </button>
  <button class="btn btn--icon" on:click={e => pull(e, $currentRemote.Name)} disabled={$currentRemote.Behind == 0 || !$currentRemote.Fetch}>
    <span class="icon">{@html octicons['arrow-down'].toSVG({ "width": 14 })}</span>
    Pull
  </button>
</div>
