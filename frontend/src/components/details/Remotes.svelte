<script lang="ts">
  import octicons from "@primer/octicons";
  import addRemote from "actions/add-remote";
  import { parseResponse } from "scripts/parse-response";
  import { remoteData, remotes } from "stores/remotes";
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
    icon.innerHTML = octicons['check'].toSVG({ "width": 16 });
    icon.classList.add('success');
    setTimeout(() => {
      icon.classList.remove('success');
      icon.innerHTML = normalSvg;
      if (reenable) {
        el.removeAttribute('disabled');
      }
    }, 1000);
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
