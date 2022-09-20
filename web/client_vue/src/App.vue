<template>
	<div id="conversation">
		<div class="conversation-container">
			<div v-for="message in messages" :key="message.index">
        <div v-if="message.Broad" class="bubbleleft">
					{{ message.Data }}
				</div>
        <div v-else class="bubbleright">
					{{ message.Data }}
        </div>
      </div>
			</div>
		</div>
		<div class="input-container">
			<input @keyup.enter="sendMessage" v-model="messageText" placeholder="Enter your message">
			<button @click="sendMessage">Send message</button>
		</div>
</template>
<script>
export default {
	data() {
		return {
			messages: [],
			messageText: "",
		}
	},
	mounted() {
    this.socket = new WebSocket("ws://0.0.0.0:9100/socket") //change to public IP 
    this.socket.onmessage = (newMessage) => {
      let msg = {
        "Broad": true,
        "Data": newMessage.data
      }
      this.messages.push(msg)
      console.log(msg)
    }
	},
	methods: {
    sendMessage(){
      let msg = {
        "Broad": false,
        "Data": this.messageText
      }
      this.socket.send(this.messageText)
      this.messages.push(msg)
      console.log(JSON.stringify(msg))
      this.messageText = ""
    }
	}
}
</script>
<style scoped>
.conversation-container {
  margin: 0 auto;
  max-width: 400px;
  height: 600px;
  padding: 0 20px;
  border: 3px solid #f1f1f1;
  overflow: scroll;
}
.input-container {
  margin: 0 auto;
  max-width: 400px;
  height: 800px;
}
.bubble-container {
	text-align: left;
}
.bubbleright {
  border: 2px solid #f1f1f1;
  background-color: #7fbefc;
  border-radius: 5px;
  padding: 10px;
  margin: 10px 0;
	width: 230px;
	float: right;
}
.bubbleleft {
  border: 2px solid #f1f1f1;
  background-color: #fafdfd;
  border-radius: 5px;
  padding: 10px;
  margin: 10px 0;
	width: 230px;
	float: left;
}
.bubble {
	background-color: #abf1ea;
	border: 2px solid #87E0D7;
	float: left;
}
.name {
  padding-right: 8px;
	font-size: 11px;
}
::-webkit-scrollbar {
  width: 10px;
}
::-webkit-scrollbar-track {
  background: #f1f1f1;
}
::-webkit-scrollbar-thumb {
  background: #888;
}
::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>