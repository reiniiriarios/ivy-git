<script lang="ts">
  import octicons from "@primer/octicons";
  import addRemote from "actions/remote/add";
  import fetchRemote from "actions/remote/fetch";
  import pullRemote from "actions/remote/pull";
  import pushRemote from "actions/remote/push";
  import { remotes } from "stores/remotes";

  function pull(e: MouseEvent | KeyboardEvent, remote: string) {
    pullRemote(remote, e.target as HTMLElement);
  }

  function push(e: MouseEvent | KeyboardEvent, remote: string) {
    pushRemote(remote, e.target as HTMLElement);
  }

  function fetch(e: MouseEvent | KeyboardEvent, remote: string) {
    fetchRemote(remote, e.target as HTMLElement);
  }

  const newRemote = () => addRemote();
</script>

<div class="remotes">
  <div class="remotes__header">
    <h2>Remotes</h2>
    <div>
      <button class="btn btn-sm" on:click={newRemote}>
        Add
        {@html octicons.plus.toSVG({width: 14})}
      </button>
    </div>
  </div>
    {#if Object.entries($remotes).length}
    <table>
      <thead>
        <tr>
          <th>Remote</th>
          <th>Repo</th>
          <th>Ahead/Behind</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each Object.entries($remotes) as [_, remote]}
          <tr class="remote" data-name="{remote.Name}" data-menu="remote">
            <td>
              <div class="remote__name">{remote.Name}</div>
              <div class="remote__fetch-push">
                {#if remote.Fetch && remote.Push}
                  {@html octicons['arrow-switch'].toSVG({ "width": 16 })}
                {:else if remote.Fetch}
                  {@html octicons['arrow-left'].toSVG({ "width": 16 })}
                {:else if remote.Push}
                  {@html octicons['arrow-right'].toSVG({ "width": 16 })}
                {:else}
                  {@html octicons['dash'].toSVG({ "width": 16 })}
                {/if}
              </div>
            </td>
            <td>
              <div class="remote__site">{remote.Site}</div>
              <div class="remote__repo"><a href="{remote.Url}" target="_blank">{remote.Repo}</a></div>
            </td>
            <td>
              <div class="remote__ahead">
                {@html octicons['arrow-up'].toSVG({ "width": 16 })}
                {remote.Ahead}
              </div>
              <div class="remote__behind">
                {@html octicons['arrow-down'].toSVG({ "width": 16 })}
                {remote.Behind}
              </div>
            </td>
            <td class="remote__actions">
              <button class="btn" on:click={e => fetch(e, remote.Name)}>
                <span class="icon">{@html octicons['arrow-down-left'].toSVG({ "width": 16 })}</span>
                Fetch
              </button>
              <button class="btn" on:click={e => push(e, remote.Name)} disabled={remote.Ahead == 0 || !remote.Push}>
                <span class="icon">{@html octicons['arrow-up'].toSVG({ "width": 16 })}</span>
                Push
              </button>
              <button class="btn" on:click={e => pull(e, remote.Name)} disabled={remote.Behind == 0 || !remote.Fetch}>
                <span class="icon">{@html octicons['arrow-down'].toSVG({ "width": 16 })}</span>
                Pull
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <div class="remotes__no-remote">No remotes.</div>
  {/if}
</div>
