<script>
  import { onMount } from "svelte";
  import { dialog } from "../stores/dialog.js";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";

  onMount(() => {});

  function close() {
    $state = "nothing";
  }
</script>

<div id="dialog">
  <span>{$dialog}</span>
  <div id="buttonbar">
    <button
      id="closebtn"
      style="background-color: {$theme.textAreaColor}; color: {$theme.textColor};"
      on:click={close}
    >
      close
    </button>
  </div>
</div>

<style>
  #dialog {
    display: flex;
    flex-direction: column;
    padding: 0px;
    margin: 0px;
  }

  #buttonbar {
    display: flex;
    flex-direction: row;
    margin: 10px;
    padding: 10px;
    align-items: center;
    justify-content: center;
  }

  button {
    padding: 5px;
    border-radius: 10px;
  }
</style>
