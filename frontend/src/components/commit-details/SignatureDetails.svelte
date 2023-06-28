<script lang="ts">
  import octicons from "@primer/octicons";
  import { commitSignature } from "stores/commit-details";
</script>

<div class="signature-details">
  {#if $commitSignature?.Status}
    {#if $commitSignature.Status !== 'N'}
      <span class="gpg-status gpg-status--{$commitSignature.Status}">
        <span class="gpg-status__icon">{@html octicons.verified.toSVG({width: 16})}</span>
        <span class="gpg-status__tip">
          <span class="gpg-status__text">
            {$commitSignature.Description}
          </span>
        </span>
      </span>
      <span class="commit-details__gpg-key" data-menu="text">{$commitSignature.Key}</span>
      {#if $commitSignature.Name}
        <span class="commit-details__gpg-name" data-menu="text">{$commitSignature.Name}</span>
      {/if}
    {:else}
      None
    {/if}
  {/if}
</div>
