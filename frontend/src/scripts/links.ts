import { BrowserOpenURL } from "wailsjs/runtime/runtime";

export function addLinkListener() {
  document.body.addEventListener("click", function (e: MouseEvent) {
    if (e.target instanceof HTMLAnchorElement) {
      e.preventDefault();
      console.log("Capturing link:", e.target.innerText);
      BrowserOpenURL(e.target.href);
    }
  });
}
