import checkIcon from "scripts/button-check-icon";
import { parseResponse } from "scripts/parse-response";
import { remoteData } from "stores/remotes";
import { FetchRemote } from "wailsjs/go/main/App";

function fetchRemote(remote: string, button?: HTMLElement) {
  if (button) button.setAttribute('disabled', 'disabled');
  FetchRemote(remote).then((r) => {
    parseResponse(r, () => {
      if (button) checkIcon(button, true);
      remoteData.refresh();
    }, () => {
      if (button) button.removeAttribute('disabled')
    });
  });
}

export default fetchRemote;
