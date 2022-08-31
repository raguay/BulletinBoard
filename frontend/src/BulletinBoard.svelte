<script>
  import { onMount, afterUpdate } from "svelte";
  import Message from "./components/Message.svelte";
  import Dialog from "./components/Dialog.svelte";
  import Raw from "./components/Raw.svelte";
  import { state } from "./stores/state.js";
  import { theme } from "./stores/theme.js";
  import { message } from "./stores/message.js";
  import { dialog } from "./stores/dialog.js";
  import { raw } from "./stores/raw.js";
  import * as rt from "../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  let containerDOM = null;
  let minWidth = 300;
  let minHeight = 60;

  onMount(async () => {
    $state = "nothing";
    await getTheme();

    //
    // Set a function to run when a event (signal) is sent from the webserver.
    //
    rt.EventsOn("message", (msg) => {
      if (msg.trim().length !== 0) {
        //
        // Set the message state and save the message in the store.
        //
        $state = "message";
        $message = msg;
      } else {
        //
        // An empty message send by having just a space turns off the BulletinBoard.
        //
        $state = "nothing";
      }
    });
  });

  afterUpdate(() => {
    //
    // The nothing state should force a window hiding. Otherwise, show the window.
    //
    if ($state === "nothing") {
      rt.WindowHide();
    } else {
      rt.WindowShow();
    }

    //
    // Figure out the width and height of the new canvas.
    //
    if (containerDOM !== null) {
      let width = minWidth;
      let height = minHeight;
      if (height < containerDOM.clientHeight)
        height = containerDOM.clientHeight;
      if (width < containerDOM.clientWidth) width = containerDOM.clientWidth;
      rt.WindowSetSize(width, height);
    }
  });

  async function getTheme(callback) {
    //
    // This would read the theme from a file. It currently just sets a typical theme.
    // I love the Dracula color theme.
    //
    $theme = {
      font: "Fira Code, Menlo",
      fontSize: "12pt",
      textAreaColor: "#454158",
      backgroundColor: "#22212C",
      textColor: "#80ffea",
      borderColor: "#1B1A23",
      Cyan: "#80FFEA",
      Green: "#8AFF80",
      Orange: "#FFCA80",
      Pink: "#FF80BF",
      Purple: "#9580FF",
      Red: "#FF9580",
      Yellow: "#FFFF80",
    };
    if (typeof callback !== "undefined") callback();
  }

  function wait(ms) {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(ms);
      }, ms);
    });
  }
</script>

<div
  id="closure"
  bind:this={containerDOM}
  style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; font-family: {$theme.font}; font-size: {$theme.fontSize};"
>
  <div id="header" data-wails-drag>
    <h3>Bulletin Board</h3>
  </div>
  <div id="main">
    {#if $state === "message"}
      <Message />
    {:else if $state === "dialog"}
      <Dialog />
    {:else if $state === "raw"}
      <Raw />
    {/if}
  </div>
</div>

<style>
  :global(body) {
    margin: 0px;
    padding: 0px;
    overflow: hidden;
    border: transparent solid 1px;
    border-radius: 10px;
    background-color: transparent;
  }

  #closure {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 0px;
    border-radius: 10px;
    overflow: hidden;
  }

  #header {
    height: 20px;
    margin: 0px;
    padding: 5px;
    -webkit-user-select: none;
    user-select: none;
    cursor: default;
  }

  #main {
    display: flex;
    flex-direction: column;
    margin: 0px 0px 0px 20px;
    padding: 0px;
    min-width: 100px;
  }

  h3 {
    text-align: center;
    margin: 0px;
    padding: 0px;
    cursor: default;
    font-size: 1 em;
  }
</style>
