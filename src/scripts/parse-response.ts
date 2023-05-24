import { messageDialog } from "stores/message-dialog";

interface DataResponse {
  Response: string;
  Message: string;
  [data: string]: any;
}

type ResponseCallback = () => void | any;

export function parseResponse(response: DataResponse, success?: ResponseCallback, failure?: ResponseCallback) {
  if (response.Response === "error") {
    messageDialog.error({
      message: response.Message
    });
    if (failure) failure();
    return
  }
  if (success) success();
}
