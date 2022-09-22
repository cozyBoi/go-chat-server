window.addEventListener("load", function(evt) {
    var loc = window.location;
    var uri = 'http:';
    var roomNumber;

    if (loc.protocol === 'https:') {
       uri = 'https:';
    }
    uri += '//' + loc.host;
    uri += loc.pathname;
    console.log(uri)

    //1. get room number
    httpGetAsync(loc.pathname + 'rooms', printRooms)

    document.getElementById("create").onclick = function(evt) {
        //=> post chat room
        console.log("asdsd")
    };

    function enterRoom(evt) {
        console.log("asdsd")
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
            document.location.href = uri + "rooms/1";
        }
        xmlHttp.open("GET", uri + "rooms/1", true); // true for asynchronous 
        xmlHttp.send(null);
    };

    function printRooms(room){
        //TODO add div
        console.log(room);
        var parseRooms = room.split(",");
        var i;
        for(i = 0; i < parseRooms.length; i++){
            let btn = document.createElement("button");
            btn.innerHTML = parseRooms[i];
            btn.onclick = enterRoom;
            document.body.appendChild(btn);
        }
        //var newDiv = document.createElement("div");
        //newDiv.appendChild(btn);
    }

    function httpGetAsync(theUrl, callback)
    {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(xmlHttp.responseText);
        }
        xmlHttp.open("GET", theUrl, true); // true for asynchronous 
        xmlHttp.send(null);
    }
});