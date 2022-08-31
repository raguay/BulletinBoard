<script>
  import { message } from "../stores/message.js";
</script>

<div id="message">
  <span>{$message}</span>
</div>

<style>
  #message {
    display: flex;
    flex-direction: column;
    margin: 0px;
    padding: 10px;
  }
</style>
