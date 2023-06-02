import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { inProgressCommitMessage } from "stores/ui";
import { RebaseOnBranch } from "wailsjs/go/main/App";

function rebaseOnBranch(branch: string) {
  RebaseOnBranch(branch).then(r => {
    parseResponse(r, () => {
      commitData.refresh();
      commitSignData.refresh();
    }, () => {
      inProgressCommitMessage.fetch();
    });
  });
}

export default rebaseOnBranch;
