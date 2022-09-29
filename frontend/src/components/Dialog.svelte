<script>
  import { afterUpdate, onMount } from "svelte";
  import { dialog } from "../stores/dialog.js";
  import { state } from "../stores/state.js";
  import { theme } from "../stores/theme.js";
  import * as rt from "../../wailsjs/runtime/runtime.js";

  let style = `background-color: ${$theme.backgroundColor}; color: ${$theme.textColor}; border-color: ${$theme.borderColor};`;
  let buttonStyle = `background-color: ${$theme.backgroundColor}; color: ${$theme.textColor}; border-color: ${$theme.borderColor}; box-shadow: ${$theme.boxShadow};`;
  let radiogroup;

  let inputTypes = [
    "input",
    "selection",
    "radio",
    "checkbox",
    "color",
    "date",
    "datetime",
    "email",
    "file",
    "month",
    "password",
    "tel",
    "time",
    "url",
    "week",
  ];

  onMount(() => {});

  afterUpdate(() => {
    //
    // Make the first input type element be focused.
    //
    let elem = window.document.getElementById(
      $dialog.items.filter((item) => inputTypes.includes(item.modaltype))[0].id
    );
    elem.focus();
    rt.WindowCenter();
  });

  function buttonClick(action) {
    let oneradio = false;
    switch (action) {
      case "submit":
        rt.EventsEmit(
          "dialogreturn",
          $dialog.items
            .filter((item) => {
              let ret = inputTypes.includes(item.modaltype);
              if (item.modaltype === "radio") {
                if (!oneradio) {
                  oneradio = true;
                  return true;
                } else {
                  return false;
                }
              } else {
                return ret;
              }
            })
            .map((item) => {
              switch (item.modaltype) {
                case "radio":
                  return { name: item.name, value: radiogroup };

                default:
                  return { name: item.name, value: item.value };
              }
            })
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

  function processKey(e) {
    //
    // If an Enter key is pressed, run the submit process.
    //
    if (e.key === "Enter") buttonClick("submit");
  }
</script>

<div id="dialogOuter" on:keydown={processKey}>
  <!-- Figure the items to display -->
  {#each $dialog.items as item}
    {#if item.modaltype === "label"}
      <label id={item.id} name={item.name} {style} for={item.forid}>
        {item.value}
      </label>
    {:else if item.modaltype === "input"}
      <input
        id={item.id}
        name={item.name}
        {style}
        type="text"
        bind:value={item.value}
      />
    {:else if item.modaltype === "selection"}
      <select id={item.id} name={item.name} bind:value={item.value} {style}>
        {#each $dialog.items as options}
          {#if options.modaltype === "option"}
            <option value={options.value} {style}>{options.value}</option>
          {/if}
        {/each}
      </select>
    {:else if item.modaltype === "radio"}
      <div class="horzdiv">
        <input
          type="radio"
          id={item.id}
          {style}
          name={item.name}
          bind:group={radiogroup}
          value={item.value}
        />
        <label for={item.name}>{item.value}</label>
      </div>
    {:else if item.modaltype === "checkbox"}
      <div class="horzdiv">
        <input
          type="checkbox"
          id={item.id}
          {style}
          name={item.name}
          bind:checked={item.value}
        />
        <label for={item.name}>{item.for}</label>
      </div>
    {:else if item.modaltype === "color"}
      <input
        type="color"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "date"}
      <input
        type="date"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "datetime"}
      <input
        type="datetime-local"
        {style}
        id={item.id}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "email"}
      <input
        type="email"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "file"}
      <input
        type="file"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "month"}
      <input
        type="month"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "password"}
      <input
        type="password"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "tel"}
      <input
        type="tel"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "time"}
      <input
        type="time"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "url"}
      <input
        type="time"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {:else if item.modaltype === "week"}
      <input
        type="week"
        id={item.id}
        {style}
        name={item.name}
        bind:value={item.value}
      />
    {/if}
  {/each}
  <div id="buttonbar">
    <!-- Set the buttons called for -->
    {#each $dialog.buttons as button}
      <button
        id={button.id}
        name={button.name}
        style={buttonStyle}
        on:click={() => buttonClick(button.action)}
      >
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
    margin: 10px 0px 0px 0px;
  }

  #buttonbar {
    display: flex;
    flex-direction: row;
    margin: 0px;
    padding: 10px 0px 0px 0px;
    align-items: center;
    justify-content: center;
  }

  button {
    padding: 7px;
    border-radius: 10px;
    margin: 0px 20px 0px 0px;
  }

  input,
  input:active,
  select,
  select:active {
    outline-style: none;
    margin: 5px 10px;
  }

  label {
    user-select: none;
    -webkit-user-select: none;
  }

  .horzdiv {
    display: flex;
    flex-direction: row;
    margin: 0px;
    padding: 0px;
  }
</style>
