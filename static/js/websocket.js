define(function (){

    var my = {},
        conn;

    my.connect = function(commands) {
        if (window["WebSocket"]) {
            conn = new WebSocket("wss://" + window.location.host +"/ws");
            // connection closed
            conn.onclose = function (evt) {
                console.log("Connection closed.");
            }
            // recieving data
            conn.onmessage = function (evt) {
                commands.push(evt.data);
            }
        } else {
            alert("Your browser does not support WebSockets.");
        }
        return conn;
    }

    // Make the function wait until the connection is made...
    my.afterConnect = function(socket, callback) {
        setTimeout(function () {
            if (socket.readyState === 1) {
                if (callback != null) {
                    callback();
                }
                return;
            } else {
                my.afterConnect(socket, callback);
            }
        }, 5);
    }

    return my;
});