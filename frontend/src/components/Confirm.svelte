<script lang="ts">
  let cb: Function = () => {};

  (window as any).confirmModal = (question: string, callback: Function, yes: string = 'Yes', no: string = 'No') => {
    let c = document.getElementById("modal-confirm");
    let q = document.getElementById("modal-confirm-question");
    let y = document.getElementById("modal-confirm-yes");
    let n = document.getElementById("modal-confirm-no");
    q.innerText = question;
    y.innerText = yes;
    n.innerText = no;
    c.style.display = "block";
    cb = callback;
  };

  function close() {
    document.getElementById("modal-confirm").style.display = "none";
  }

  function yes() {
    cb();
    close();
  }
</script>

<div id="modal-confirm">
  <div class="overlay">
    <div class="modal">
      <div id="modal-confirm-question">...</div>
      <div class="answer">
        <button id="modal-confirm-yes" class="btn yes" on:click={yes} on:keyup={yes}>Yes</button>
        <button id="modal-confirm-no" class="btn no" on:click={close} on:keyup={close}>No</button>
      </div>
    </div>
  </div>
</div>

<style lang="scss">
  #modal-confirm {
    display: none;

    .overlay {
      z-index: 100;
    }
  }

  #modal-confirm-question {
    padding: 1rem;
  }

  .modal {
    z-index: 200;
    background-color: var(--color-modal-bg);
    display: flex;
    justify-content: space-between;
    flex-direction: column;

    .answer {
      border-top: 1px solid var(--color-modal-border);
      padding: 1rem;
    }
  }
</style>
