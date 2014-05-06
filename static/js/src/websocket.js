/* jshint undef: true, unused: true, strict: true */
/* global define, console, window, document, alert */
define(function () {

    "use strict";

    var my = {},
        conn;

    /**
     * float64 - server unix.nano
     */
    my.timestamp = 0;

    /**
     *
     * @param callback - onConnect
     * @param callback - onMessage
     */
    my.connect = function (onConnect, onMessage) {
        if (window.WebSocket) {
            // Let's open a web socket to the server
            conn = new window.WebSocket("ws://" + document.location.host + "/ws/");
            conn.binaryType = "arraybuffer";
            conn.onopen = function (data) {
                console.log("connection was opened to '" + data.currentTarget.URL + '"');
                onConnect();
            };

            conn.onerror = function () {
                console.log("connection error");
            };

            conn.onmessage = function (evt) {
                var buf = new DataStream(evt.data);
                my.timestamp = buf.readFloat64();
                var msgType = buf.readUint8();
                // world state update
                if(msgType === 1) {

                }
                // respond to a ping request
                if(msgType === 2) {
                    my.send(my.newMessage(2));
                }
                //onMessage(buf);
            };
            conn.onclose = function (data) {
                // var code = event.code;
                // var reason = event.reason;
                // var wasClean = event.wasClean;
                console.log("connection was closed to '" + data.currentTarget.URL + '"');
            };
        } else {
            alert("Your browser does not support WebSockets. :'|");
        }
    };

    my.close = function () {
        conn.close();
    };

    my.newMessage = function(msgType) {
        var cmd = new DataStream();
        cmd.writeFloat64(my.timestamp, DataStream.LITTLE_ENDIAN);
        cmd.writeUint8(msgType);
        return cmd
    }

    /**
     *
     * @param message
     * @returns {boolean}
     */
    my.send = function (message) {
        // Wait until previous message has been sent
        // while(conn.bufferedAmount === 0) {}
        if (conn.readyState != window.WebSocket.OPEN) {
            console.log('Connection is not ready');
            return false;
        }
        conn.send(message.buffer);
        return true;
    };

    return my;
});