<script>
  import { dialog } from "../stores/dialog.js";
  import { state } from "../stores/state.js";
  import * as rt from "../../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  function buttonClick(action) {
    switch (action) {
      case "Submit":
        rt.EventsEmit(
          "dialogreturn",
          $dialog.items.filter((item) => item.modaltype === "input")
        );
        $state = "nothing";
        break;

      default:
        rt.EventsEmit("dialogreturn", {
          canceled: true,
        });
        $state = "nothing";
        break;
    }
  }
</script>

<div id="dialogOuter">
  <!-- Figure the items to display -->
  {#each $dialog.items as item}
    {#if item.modaltype === "label"}
      <label id={item.id} for={item.for}>{item.value}</label>
    {:else if item.modaltype === "input"}
      <input id={item.id} type="text" bind:value={item.value} />
    {/if}
  {/each}
  <div id="buttonbar">
    <!-- Set the buttons called for -->
    {#each $dialog.buttons as button}
      <button id={button.id} on:click={() => buttonClick(button.action)}>
        {button.name}
      </button>
    {/each}
  </div>
</div>

<style>
  #dialogOuter {
    display: flex;
    flex-direction: column;
    padding: 0px;
    margin: 0px;
  }

  #buttonbar {
    display: flex;
    flex-direction: row;
    margin: 0px;
    padding: 10px;
    align-items: center;
    justify-content: center;
  }

  button {
    padding: 5px;
    border-radius: 10px;
    margin: 0px 20px 0px 0px;
  }
</style>
