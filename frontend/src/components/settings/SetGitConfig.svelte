<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { gitConfig } from "stores/git-config";
</script>

<div class="setting setting--split">
  <div class="setting__split setting__split--left">
    <h4 class="setting__name">Global Settings</h4>
    <div class="override" class:overridden={$gitConfig.local.UserName}>
      <TextInput bind:value={$gitConfig.global.UserName} display="User Name" update={value => gitConfig.setUserName('global', value)} />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserEmail}>
      <TextInput bind:value={$gitConfig.global.UserEmail} display="User Email" update={value => gitConfig.setUserEmail('global', value)} />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserSigningKey}>
      <TextInput bind:value={$gitConfig.global.UserSigningKey} display="User Signing Key" update={value => gitConfig.setUserSigningKey('global', value)} />
    </div>
    <div class="override" class:overridden={$gitConfig.local.CommitGpgSign}>
      <Checkbox bind:checked={$gitConfig.global.CommitGpgSign} display="Sign Commits" on:click={() => gitConfig.setSignCommits('global', !$gitConfig.local.CommitGpgSign)} />
    </div>
  </div>
  <div class="setting__split setting__split--right">
    <h4 class="setting__name">Local (Repo) Settings</h4>
    <div class="override" class:empty={!$gitConfig.local.UserName}>
      <TextInput bind:value={$gitConfig.local.UserName} display="User Name" update={value => gitConfig.setUserName('local', value)} />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserEmail}>
      <TextInput bind:value={$gitConfig.local.UserEmail} display="User Email" update={value => gitConfig.setUserEmail('local', value)} />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserSigningKey}>
      <TextInput bind:value={$gitConfig.local.UserSigningKey} display="User Signing Key" update={value => gitConfig.setUserSigningKey('local', value)} />
    </div>
    <div class="override" class:empty={!$gitConfig.local.CommitGpgSign}>
      <Checkbox bind:checked={$gitConfig.local.CommitGpgSign} display="Sign Commits" on:click={() => gitConfig.setSignCommits('local', !$gitConfig.local.CommitGpgSign)} />
    </div>
  </div>
</div>
