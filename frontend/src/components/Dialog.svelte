<script>
  import { dialog } from "../stores/dialog.js";
</script>

<div id="dialogOuter">
  {@html $dialog.html}
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
