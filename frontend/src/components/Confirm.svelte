<script lang="ts">
  let cb: Function = () => {};

  (window as any).confirmModal = (question: string, callback: Function, yes: string = 'Yes', no: string = 'No') => {
    let c = document.getElementById("modal-confirm");
    let q = document.getElementById("modal-confirm-text");
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

<div class="modal" id="modal-confirm">
  <div class="overlay">
    <div class="modal__box">
      <div class="modal__text" id="modal-confirm-text">...</div>
      <div class="modal__response">
        <button id="modal-confirm-yes" class="btn yes" on:click={yes} on:keyup={yes}>Yes</button>
        <button id="modal-confirm-no" class="btn no" on:click={close} on:keyup={close}>No</button>
      </div>
    </div>
  </div>
</div>
