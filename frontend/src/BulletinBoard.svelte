<script>
  import { onMount, afterUpdate, tick } from "svelte";
  import Message from "./components/Message.svelte";
  import Dialog from "./components/Dialog.svelte";
  import Raw from "./components/Raw.svelte";
  import { state } from "./stores/state.js";
  import { theme } from "./stores/theme.js";
  import { message } from "./stores/message.js";
  import { raw } from "./stores/raw.js";
  import { dialog } from "./stores/dialog.js";
  import * as rt from "../wailsjs/runtime/runtime.js"; // the runtime for Wails2

  let containerDOM = null;
  let minWidth = 300;
  let minHeight = 60;
  let width = 300;
  let height = 60;

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
    rt.EventsOn("append", (msg) => {
      $state = "message";
      $message = $message + msg;
    });
    rt.EventsOn("dialog", (msg) => {
      $state = "raw";
      $raw = msg;
    });
    rt.EventsOn("modal", (msg) => {
      $state = "dialog";
      $dialog = msg;
    });
  });

  afterUpdate(async () => {
    //
    // Figure out the width and height of the new canvas.
    //
    await tick();
    await tick();
    await tick();
    await tick();
    if (containerDOM !== null) {
      if ($state === "raw") {
        //
        // Check the overrides for a raw state.
        //
        width = minWidth;
        height = minHeight;
        if (width < $raw.width) width = $raw.width;
        if (height < $raw.height) height = $raw.height;

        //
        // Set the position on the screen.
        //
        await rt.WindowSetSize(width, height);
        rt.WindowSetPosition($raw.y, $raw.x);
      } else {
        //
        // Non-raw dialogs have a calculated size.
        //
        width = minWidth;
        height = minHeight;
        if (height < containerDOM.clientHeight)
          height = containerDOM.clientHeight;
        if (width < containerDOM.clientWidth) width = containerDOM.clientWidth;

        //
        // Set to the determined size;
        //
        await rt.WindowSetSize(width, height);
      }
    }

    //
    // The nothing state should force a window hiding. Otherwise, show the window.
    //
    if ($state === "nothing") {
      //
      // Hide the window.
      //
      rt.WindowHide();
    } else {
      //
      // Show the window.
      //
      rt.WindowShow();
    }
  });

  async function getTheme() {
    //
    // This would read the theme from a file. It currently just sets a typical theme.
    // I love the Dracula color theme.
    //
    $theme = {
      name: "Default",
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
      boxShadow: "2px 2px 2px #9580ff90",
    };
  }
</script>

<div
  id="closure"
  style="background-color: {$theme.backgroundColor}; color: {$theme.textColor}; font-family: {$theme.font}; font-size: {$theme.fontSize};"
  bind:this={containerDOM}
>
  <div id="header">
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
    padding: 10px;
    border: solid transparent 1px;
    overflow: hidden;
  }

  #header {
    height: 20px;
    margin-top: -10px;
    padding: 5px;
    -webkit-user-select: none;
    user-select: none;
    cursor: default;
    --wails-draggable: drag;
  }

  #main {
    display: flex;
    flex-direction: column;
    margin: 0px;
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

  h3 {
    user-select: none;
    -webkit-user-select: none;
  }
</style>
