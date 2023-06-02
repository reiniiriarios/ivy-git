import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { PushTag } from "wailsjs/go/main/App";

function pushTag(tag: string) {
  PushTag(tag).then(r => {
    parseResponse(r, () => {
      commitData.refresh();
      commitSignData.refresh();
    });
  });
}

export default pushTag;
