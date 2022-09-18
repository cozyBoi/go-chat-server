<template>
  <form :action="sendMessage" @click.prevent="onSubmit">
    <input v-model="message" type="text" name="" id="">
    <input type="submit" value="Send" @click="sendMessage">
  </form>
  <div>
    <h3>Message in a WebSocket</h3>
    <h6>{{ message }}</h6>
  </div>
</template>

<script>

import Talk from "talkjs";

export default {
  name: 'App',
  data() {
    return {
      message: "",
      socket: null,
      recvMsg: "",
      chatRoom: "",
    }
  },
  mounted(){
    this.socket = new WebSocket("ws://localhost:9100/socket")
    this.socket.onmessage = (msg) => {
      this.recvMsg = msg.data
      //this.chatRoom.append(this.recvMsg)
    }
  },
  methods: {
    sendMessage(){
      let msg = {
        "greeting": this.message
      }
      this.socket.send(JSON.stringify(msg))
    }
  },
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
