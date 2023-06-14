import { parseResponse } from "scripts/parse-response";
import { commitData, commitSignData } from "stores/commits";
import { messageDialog } from "stores/message-dialog";
import { MergeFastForward } from "wailsjs/go/main/App";

function fastForwardMerge(branch: string) {
  messageDialog.confirm({
    heading: 'Fast-forward Merge',
    message: `Merge the current branch into <strong>${branch}</strong> via fast-forward only?`,
    confirm: 'Merge',
    okay: 'Cancel',
    callbackConfirm: () => {
      MergeFastForward(branch).then(r => {
        parseResponse(r, () => {
          commitData.refresh();
          commitSignData.refresh();
        });
      });
    },
  });
}

export default fastForwardMerge;
