define(["datastream", "entity"], function (DataStream, entity){

    var my = {},
        conn;

    my.connect = function(afterconnect, onMessage) {
        if (window["WebSocket"]) {
            // Let us open a web socket
            conn = new WebSocket("wss://"+document.location.host+"/ws/");
            conn.binaryType = "arraybuffer";
            conn.onopen = function() {
                console.log("connection open");
                afterconnect();
            }
            conn.onerror = function() {
                console.log("connection error");
            }
            conn.onmessage = onMessage;
            conn.onclose = function() {
                console.log("Connection was closed");
            };
        } else {
            alert("Your browser does not support WebSockets. :'|");
        }
    }

    my.send = function(message) {
        if (conn.readyState != WebSocket.OPEN) {
            return false;
        }
        conn.send(message);
        return true;
    }

    // Make the function wait until the connection is made...
    return my;
});