import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commit-data";
import { PullRemoteBranch } from "wailsjs/go/main/App";

function pullBranch(branch: string, remote: string) {
  PullRemoteBranch(remote, branch, true).then(r => {
    parseResponse(r, () => {
      commitData.refresh();
      commitSignData.refresh();
    });
  });
}

export default pullBranch;
