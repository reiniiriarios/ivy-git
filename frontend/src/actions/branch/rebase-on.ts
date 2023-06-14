import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { commitMessage } from "stores/ui";
import { RebaseOnBranch } from "wailsjs/go/main/App";

function rebaseOnBranch(branch: string) {
  RebaseOnBranch(branch).then(r => {
    parseResponse(r, () => {
      commitData.refresh();
      commitSignData.refresh();
    }, () => {
      commitMessage.fetch();
    });
  });
}

export default rebaseOnBranch;
