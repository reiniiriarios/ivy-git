<script lang="ts">
  import { currentCommit } from "stores/commit-details";
  import { commits, commitsMap } from "stores/commits";
  import { currentTab } from "stores/ui";

  export let subject: string;
  export let body: string;

  function codify(s: string): string {
    if (!s) return "";
    s = s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    return s.replaceAll(/`([^`]+?)`/g, '<code>$1</code>');
  }

  // Add spans for commit hashes found, or things that look like them.
  function commitLinks(s: string): string {
    if (!s) return "";
    return s.replaceAll(/([a-f0-9]{40})/gi, m => {
      // Only add if the hash is in the commits list.
      let commit_id = $commitsMap.get(m);
      if (typeof commit_id === 'number') {
        let span = document.createElement('span');
        span.dataset.id = commit_id.toString();
        span.classList.add('linked-commit');
        span.classList.add('linked-commit--linkable');
        span.innerHTML = m;
        return span.outerHTML;
      }
      return m;
    });
  }

  // If the commit body is clicked, look to see if the click was on a commit hash. If so, clicky clicky.
  function clickBody(e: (MouseEvent | KeyboardEvent) & { currentTarget: HTMLElement }) {
    let el = e.target as HTMLElement;
    if (el.classList.contains('linked-commit')) {
      let commit_id = parseInt(el.dataset.id);
      if ($commits[commit_id]) {
        currentTab.set('tree');
        currentCommit.set($commits[commit_id]);
      }
    }
  }
</script>

<div class="message">
  <div class="message__subject">{@html codify(subject)}</div>
  <div id="commit-details-body" on:click={clickBody} on:keypress={clickBody} class="message__body">
    <!-- body is markup html, do not escape -->
    {@html commitLinks(body)}
  </div>
</div>
