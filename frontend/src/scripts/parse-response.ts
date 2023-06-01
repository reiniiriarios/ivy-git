import { messageDialog } from "stores/message-dialog";

interface DataResponse {
  Response: string;
  Message: string;
  [data: string]: any;
}

type ResponseCallback = () => void | any;

export function parseResponse(response: DataResponse, success?: ResponseCallback, failure?: ResponseCallback, noError: boolean = false) {
  if (response.Response === "error") {
    if (!noError) {
      messageDialog.error({
        message: response.Message
      });
    }
    if (failure) failure();
    return
  }
  if (success) success();
}
