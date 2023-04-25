<script lang="ts">
  let cb: Function = () => {};

  (window as any).messageModal = (message: string, callback: Function = cb) => {
    let c = document.getElementById("modal-message");
    let q = document.getElementById("modal-message-text");
    q.innerText = message;
    c.style.display = "block";
    cb = callback;
  };

  function close() {
    document.getElementById("modal-message").style.display = "none";
  }
</script>

<div class="modal" id="modal-message">
  <div class="overlay">
    <div class="modal__box">
      <div class="modal__text" id="modal-message-text">...</div>
      <div class="modal__response">
        <button class="btn okay" on:click={close} on:keyup={close}>Okay</button>
      </div>
    </div>
  </div>
</div>
