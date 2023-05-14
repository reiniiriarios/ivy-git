<script lang="ts">
  import octicons from "@primer/octicons";
  import { remotes } from "stores/remotes";
  import { onMount } from 'svelte';
  import { FetchRemote } from "wailsjs/go/main/App";

  function pull(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    FetchRemote(remote).then((r) => {
      if (r.Response === "error") {
        (window as any).messageModal(r.Message);
        el.removeAttribute('disabled');
        remotes.refresh();
      } else {
        checkIcon(el);
      }
    });
  }

  function push(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    FetchRemote(remote).then((r) => {
      if (r.Response === "error") {
        (window as any).messageModal(r.Message);
        el.removeAttribute('disabled');
        remotes.refresh();
      } else {
        checkIcon(el);
      }
    });
  }

  function fetch(e: MouseEvent | KeyboardEvent, remote: string) {
    let el = e.target as HTMLElement;
    el.setAttribute('disabled', 'disabled');
    FetchRemote(remote).then((r) => {
      if (r.Response === "error") {
        (window as any).messageModal(r.Message);
        el.removeAttribute('disabled');
        remotes.refresh();
      } else {
        checkIcon(el);
      }
    });
  }

  function checkIcon(el: HTMLElement) {
    let icon = el.getElementsByClassName('icon')[0];
    let normalSvg = icon.innerHTML;
    icon.innerHTML = octicons['check'].toSVG({ "width": 16 });
    icon.classList.add('success');
    setTimeout(() => {
      icon.classList.remove('success');
      icon.innerHTML = normalSvg;
      el.removeAttribute('disabled');
    }, 1000);
  }

  onMount(() => {
    remotes.refresh();
  })
</script>

<div class="remotes">
  <h2>Remotes</h2>
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
      {#if $remotes.length}
        {#each Object.entries($remotes) as [_, remote]}
          <tr class="remote">
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
              <button class="btn" on:click={e => pull(e, remote.Name)} disabled={remote.Behind == 0}>
                <span class="icon">{@html octicons['arrow-down'].toSVG({ "width": 16 })}</span>
                Pull
              </button>
              <button class="btn" on:click={e => push(e, remote.Name)} disabled={remote.Ahead == 0}>
                <span class="icon">{@html octicons['arrow-up'].toSVG({ "width": 16 })}</span>
                Push
              </button>
            </td>
          </tr>
        {/each}
      {:else}
        <tr class="remote">
          <td colspan="4">...</td>
        </tr>
      {/if}
    </tbody>
  </table>
  <!-- TODO: add, remove -->
</div>
