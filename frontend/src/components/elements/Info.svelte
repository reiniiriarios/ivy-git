<script lang="ts">
  import octicons from "@primer/octicons";

  export let position: string = 'right';
  export let size: number = 14;
</script>

<span class="info">
  <span class="info__icon">{@html octicons.info.toSVG({width: size})}</span>
  <span class="info__tip info__tip--{position}">
    <span class="info__text">
      <slot />
    </span>
  </span>
</span>
