<script lang="ts">
  import Info from "components/elements/Info.svelte";
  import TextInput from "components/elements/TextInput.svelte";
  import { checkGpgKeyFormat } from "scripts/check-gpg";
  import { gitConfig } from "stores/git-config";
</script>

<div class="setting setting--split">
  <div class="setting__split setting__split--left">
    <h4 class="setting__name setting__name--split">
      Global Settings
      <Info>
        Global settings overridden by local settings.
        These are typically stored in <code>~/.gitconfig</code>.
      </Info>
    </h4>
    <div class="override" class:overridden={$gitConfig.local.UserName}>
      <TextInput
        bind:value={$gitConfig.global.UserName}
        display="User Name"
        update={value => gitConfig.setUserName('global', value)}
      />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserEmail}>
      <TextInput
        bind:value={$gitConfig.global.UserEmail}
        display="User Email"
        update={value => gitConfig.setUserEmail('global', value)}
      />
    </div>
    <div class="override" class:overridden={$gitConfig.local.UserSigningKey}>
      <TextInput
        bind:value={$gitConfig.global.UserSigningKey}
        display="User Signing Key"
        validate={checkGpgKeyFormat}
        update={value => gitConfig.setUserSigningKey('global', value)}
        tip="See
          <a href='https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key' target='_blank'
          >GitHub docs</a> for more information."
      />
    </div>
  </div>
  <div class="setting__split setting__split--right">
    <h4 class="setting__name setting__name--split">
      Local (Repo) Settings
      <Info position="left">
        Local settings are per-repo and will override global settings.
        These are stored in <code>&lt;repo_dir&gt;/.git/config</code>.
      </Info>
    </h4>
    <div class="override" class:empty={!$gitConfig.local.UserName}>
      <TextInput
        bind:value={$gitConfig.local.UserName}
        display="User Name"
        update={value => gitConfig.setUserName('local', value)}
      />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserEmail}>
      <TextInput
        bind:value={$gitConfig.local.UserEmail}
        display="User Email"
        update={value => gitConfig.setUserEmail('local', value)}
      />
    </div>
    <div class="override" class:empty={!$gitConfig.local.UserSigningKey}>
      <TextInput
        bind:value={$gitConfig.local.UserSigningKey}
        display="User Signing Key"
        validate={checkGpgKeyFormat}
        update={value => gitConfig.setUserSigningKey('local', value)}
        tip="See
          <a href='https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key' target='_blank'
          >GitHub docs</a> for more information."
        tipPosition="left"
      />
    </div>
  </div>
</div>
