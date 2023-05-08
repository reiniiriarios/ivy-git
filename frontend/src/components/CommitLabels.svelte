<script lang="ts">
  import octicons from '@primer/octicons';

  import { getLabelDist, type Commit, type Ref } from "../scripts/graph";

  export let commit: Commit;
  export let HEAD: Ref;
</script>

<div class="refs">
  {#if commit.Branches && commit.Branches.length}
    {#each commit.Branches as b}
      <div class="refs__label refs__label--branch" data-name="{b.Name}">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{b.Name}</div>
        {#if commit.Remotes && commit.Remotes.length}
          {#each commit.Remotes as r}
            <div class="refs__leaf">{r.ShortName}</div>
          {/each}
        {/if}
      </div>
    {/each}
  {:else if commit.Remotes && commit.Remotes.length}
    {#each commit.Remotes as r}
      <div class="refs__label refs__label--branch" data-name="{r.Name}">
        <div class="refs__icon">{@html octicons['git-branch'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{r.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Hash == HEAD.Hash}
    <div class="refs__label refs__label--head">
      <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
      <div class="refs__label-name">HEAD</div>
      {#if commit.Heads && commit.Heads.length}
        {#each commit.Heads as h}
          <div class="refs__leaf">{h.ShortName}</div>
        {/each}
      {/if}
    </div>
  {:else if commit.Heads && commit.Heads.length}
    {#each commit.Heads as h}
      <div class="refs__label refs__label--head">
        <div class="refs__icon">{@html octicons['arrow-right'].toSVG({ "width": 14 })}</div>
        <div class="refs__leaf">{h.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Tags && commit.Tags.length}
    {#each commit.Tags as t}
      <div class="refs__label refs__label--tag" data-name="{t.Name}">
        <div class="refs__icon">{@html octicons['tag'].toSVG({ "width": 14 })}</div>
        <div class="refs__label-name">{t.Name}</div>
      </div>
    {/each}
  {/if}

  {#if commit.Stash}
    <div class="refs__label refs__label--stash">
      <div class="refs__icon">{@html octicons['inbox'].toSVG({ "width": 14 })}</div>
      <div class="refs__label-name">{commit.RefName}</div>
    </div>
  {/if}

  <div class="refs__line" style="width:{getLabelDist(commit.X)}px; right:-{getLabelDist(commit.X)}px"></div>
</div>

<style lang="scss">
  .refs {
    display: flex;
    justify-content: right;
    align-items: center;
    position: relative;
    cursor: default !important;
    user-select: none;
    -webkit-user-select: none;
    white-space: nowrap;

    &__label {
      display: inline-flex;
      justify-content: left;
      align-items: center;
      border-left: 2px solid red;
      background-color: rgba(255 0 0 / 25%);
      position: relative;
      margin-right: 0.5rem;

      &::after {
        content: '';
        height: 1px;
        width: 0.5rem;
        position: absolute;
        right: -0.5rem;
        background-color: red;
      }

      &:last-child::after {
        width: 1rem;
        right: -1rem;
      }

      &:hover {
        filter: brightness(125%);
      }

      &-name {
        padding: 0.15rem 0.5rem 0.25rem 0.5rem;
        border-left: 1px solid;
        border-color: inherit;
      }
    }

    &__line {
      height: 1px;
      position: absolute;
      width: 1rem;
      right: -1rem;
      background-color: red;
    }

    &__leaf {
      padding: 0.15rem 0.5rem 0.25rem 0.5rem;
      background-color: rgba(0 0 0 / 25%);
    }

    &__icon {
      width: 1.7rem;
      height: 1.5rem;
      padding-left: 0.1rem;
      padding-right: 0.1rem;
      display: flex;
      justify-content: center;
      align-items: center;
      fill: var(--color-text);
    }
  }
</style>
