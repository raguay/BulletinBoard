<script>
  import { onMount, afterUpdate } from "svelte";
  import { raw } from "../stores/raw.js";
  import { theme } from "../stores/theme.js";
  import { state } from "../stores/state.js";

  onMount(() => {});

  function close() {
    $state = "nothing";
  }

  afterUpdate(() => {
    insertAndExecute("dialog", $raw);
  });

  //
  // The following code comes from https://stackoverflow.com/questions/2592092/executing-script-elements-inserted-with-innerhtml
  //
  function insertAndExecute(id, text) {
    let domelement = document.getElementById(id);
    domelement.innerHTML = text;
    var scripts = [];
    let ret = domelement.childNodes;
    for (var i = 0; ret[i]; i++) {
      if (
        scripts &&
        nodeName(ret[i], "script") &&
        (!ret[i].type || ret[i].type.toLowerCase() === "text/javascript")
      ) {
        scripts.push(
          ret[i].parentNode ? ret[i].parentNode.removeChild(ret[i]) : ret[i]
        );
      }
    }
    for (let script in scripts) {
      evalScript(scripts[script]);
    }
  }

  function nodeName(elem, name) {
    return elem.nodeName && elem.nodeName.toUpperCase() === name.toUpperCase();
  }

  function evalScript(elem) {
    let data = elem.text || elem.textContent || elem.innerHTML || "";
    var head =
        document.getElementsByTagName("head")[0] || document.documentElement,
      script = document.createElement("script");
    script.type = "text/javascript";
    script.appendChild(document.createTextNode(data));
    head.insertBefore(script, head.firstChild);
    head.removeChild(script);
    if (elem.parentNode) {
      elem.parentNode.removeChild(elem);
    }
  }

  //
  // End of code taken from a StackOverflow question.
  //
</script>

<div id="raw">
  <div id="dialog">
    {{ $raw }}
  </div>
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
  #raw {
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
