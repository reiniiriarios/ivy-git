<script lang="ts">
  import { linkPreviewHref } from "stores/ui";

  let hovering: boolean;
  document.addEventListener('mousemove', function (e: MouseEvent & { target: HTMLElement }) {
    if (e.target instanceof HTMLAnchorElement || e.target.parentElement instanceof HTMLAnchorElement) {
      let el = e.target instanceof HTMLAnchorElement ? e.target : e.target.parentElement as HTMLAnchorElement;
      if (!hovering) {
        linkPreviewHref.set(el.href);
        hovering = true;
      }
    } else {
      if (hovering) {
        linkPreviewHref.set('');
        hovering = false;
      }
    }
  });
</script>

<div
  class="link-preview"
  class:link-preview--visible={!!$linkPreviewHref}
>
  {$linkPreviewHref}
</div>
