<script lang="ts">
  import Checkbox from "components/elements/Checkbox.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { gitConfig } from "stores/git-config";

</script>

<div class="setting setting--split">
  <div class="setting__split setting__split--left">
    <h4 class="setting__name">Global Settings</h4>
    <div class="override" class:overridden={$gitConfig.local.UserName}>
      <TextInput bind:value={$gitConfig.global.UserName} display="User Name" />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserEmail}>
      <TextInput bind:value={$gitConfig.global.UserEmail} display="User Email" />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserSigningKey}>
      <TextInput bind:value={$gitConfig.global.UserSigningKey} display="User Signing Key" />
    </div>
    <div class="override" class:overridden={$gitConfig.local.CommitGpgSign}>
      <Checkbox bind:checked={$gitConfig.global.CommitGpgSign} display="Sign Commits" />
    </div>
  </div>
  <div class="setting__split setting__split--right">
    <h4 class="setting__name">Local Settings</h4>
    <div class="override" class:empty={!$gitConfig.local.UserName}>
      <TextInput bind:value={$gitConfig.local.UserName} display="User Name" />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserEmail}>
      <TextInput bind:value={$gitConfig.local.UserEmail} display="User Email" />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserSigningKey}>
      <TextInput bind:value={$gitConfig.local.UserSigningKey} display="User Signing Key" />
    </div>
    <div class="override" class:empty={!$gitConfig.local.CommitGpgSign}>
      <Checkbox bind:checked={$gitConfig.local.CommitGpgSign} display="Sign Commits" />
    </div>
  </div>
</div>
