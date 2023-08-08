import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { PushTag } from "wailsjs/go/ivy/App";

function pushTag(tag: string) {
  PushTag(tag).then(r => {
    parseResponse(r, () => {
      commitData.refresh();
      commitSignData.refresh();
    });
  });
}

export default pushTag;
