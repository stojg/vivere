/* jshint undef: true, unused: true, strict: true */
/* global define, console, window, document, alert */
define(function () {

    "use strict";

    var server = {},
        conn;

    /**
     * float64 - server unix.nano
     */
    server.timestamp = 0;

    /**
     *
     * @param callback - onConnect
     * @param callback - onMessage
     */
    server.connect = function (onConnect, onMessage) {
        if (!window.WebSocket) {
            alert("Your browser does not support WebSockets. :'|");
            return;
        }

        // Let's open a web socket to the server
        var proto = "ws";
        if(window.location.protocol == 'https:') {
            proto = "wss";
        }

        conn = new window.WebSocket(proto+"://" + document.location.host + "/ws/");
        conn.binaryType = "arraybuffer";
        conn.onopen = function (data) {
            console.log("connection was opened to '" + data.currentTarget.url + '"');
            onConnect();
        };

        conn.onerror = function () {
            console.log("connection error");
        };

        conn.onmessage = function (evt) {
            var buf = new DataStream(evt.data);
            server.timestamp = buf.readFloat64();
            onMessage(buf);
        };
        conn.onclose = function (event) {
            console.log("connection was closed to '" + event.currentTarget.url + '"');
        };
    };

    server.close = function () {
        conn.close();
    };

    server.newMessage = function(msgType) {
        var cmd = new DataStream();
        cmd.writeFloat64(server.timestamp, DataStream.LITTLE_ENDIAN);
        cmd.writeUint8(msgType);
        return cmd
    }

    /**
     *
     * @param message
     * @returns {boolean}
     */
    server.send = function (message) {
        // Wait until previous message has been sent
        // while(conn.bufferedAmount === 0) {}
        if (conn.readyState != window.WebSocket.OPEN) {
            console.log('Connection is not ready');
            return false;
        }
        conn.send(message.buffer);
        return true;
    };

    return server;
});
