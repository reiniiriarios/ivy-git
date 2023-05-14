interface DataResponse {
  Response: string;
  Message: string;
  [data: string]: any;
}

type ResponseCallback = () => void | any;

export function parseResponse(response: DataResponse, success?: ResponseCallback, failure?: ResponseCallback) {
  if (response.Response === "error") {
    (window as any).messageModal(response.Message);
    failure();
    return
  }
  success();
}
