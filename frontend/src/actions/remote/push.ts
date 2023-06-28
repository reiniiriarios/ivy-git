import { get } from "svelte/store";
import checkIcon from "scripts/button-check-icon";
import { parseResponse } from "scripts/parse-response";
import { remoteData } from "stores/remotes";
import { settings } from "stores/settings";
import { PushRemote } from "wailsjs/go/main/App";
import { messageDialog } from "stores/message-dialog";

function pushRemote(remote: string, button?: HTMLElement) {
  if (button) button.setAttribute('disabled', 'disabled');
  PushRemote(remote, false).then((r) => {
    if (r.Response === 'must-force' && get(settings).Workflow !== 'merge') {
      messageDialog.confirm({
        heading: 'Force Push Branch',
        message: `Unable to push, as current branch is behind its remote counterpart.\n\nForce push this branch?`,
        confirm: 'Force Push',
        okay: 'Cancel',
        callbackConfirm: () => {
          PushRemote(remote, true).then(r => {
            parseResponse(r, () => {
              if (button) checkIcon(button);
              remoteData.refresh();
            }, () => {
              if (button) button.removeAttribute('disabled');
            });
          });
        },
        callback: () => {
          if (button) button.removeAttribute('disabled');
        }
      });
    }
    else {
      parseResponse(r, () => {
        if (button) checkIcon(button);
        remoteData.refresh();
      }, () => {
        if (button) button.removeAttribute('disabled');
      });
    }
  });
}

export default pushRemote;
