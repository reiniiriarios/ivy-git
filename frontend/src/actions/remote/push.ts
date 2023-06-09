import checkIcon from "scripts/button-check-icon";
import { parseResponse } from "scripts/parse-response";
import { remoteData } from "stores/remotes";
import { PushRemote } from "wailsjs/go/main/App";

function pushRemote(remote: string, button?: HTMLElement) {
  if (button) button.setAttribute('disabled', 'disabled');
  PushRemote(remote).then((r) => {
    parseResponse(r, () => {
      if (button) checkIcon(button);
      remoteData.refresh();
    }, () => {
      if (button) button.removeAttribute('disabled')
    });
  });
}

export default pushRemote;
