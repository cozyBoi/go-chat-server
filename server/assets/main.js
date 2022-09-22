window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var printL = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        d.className = "bubbleleft"
        output.appendChild(d);  
        output.scroll(0, output.scrollHeight);
    };
    var printR = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        d.className = "bubbleright"
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        var loc = window.location;
        var uri = 'ws:';

        if (loc.protocol === 'https:') {
          uri = 'wss:';
        }
        uri += '//' + loc.host;
        uri += loc.pathname + 'ws';
        console.log(uri)

        ws = new WebSocket(uri);
        ws.onopen = function(evt) {
            //print("OPEN");
        }
        ws.onclose = function(evt) {
            //print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            printL("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            printL("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        //input.value.className = "bubble"
        printR("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});