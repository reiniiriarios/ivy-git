import { autoFetchTimer } from "events/auto-fetch";
import checkIcon from "scripts/button-check-icon";
import { parseResponse } from "scripts/parse-response";
import { remoteData } from "stores/remotes";
import { PullRemote } from "wailsjs/go/main/App";

function pullRemote(remote: string, button?: HTMLElement) {
  if (button) button.setAttribute('disabled', 'disabled');
  PullRemote(remote).then((r) => {
    parseResponse(r, () => {
      if (button) checkIcon(button);
      remoteData.refresh();
      autoFetchTimer.reset();
    }, () => {
      if (button) button.removeAttribute('disabled')
    });
  });
}

export default pullRemote;
