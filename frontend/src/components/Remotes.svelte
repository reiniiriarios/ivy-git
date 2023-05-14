<script lang="ts">
  import octicons from "@primer/octicons";
  import { remotes } from "stores/remotes";
  import { onMount } from 'svelte';

  function sync() {
    alert('todo');
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
            <button class="btn" on:click={sync}>{@html octicons['arrow-down-left'].toSVG({ "width": 16 })} Fetch</button>
            <button class="btn" on:click={sync} disabled={remote.Behind == 0}>{@html octicons['arrow-down'].toSVG({ "width": 16 })} Pull</button>
            <button class="btn" on:click={sync} disabled={remote.Ahead == 0}>{@html octicons['arrow-up'].toSVG({ "width": 16 })} Push</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
  <!-- TODO: add, remove -->
</div>
