/* jshint undef: true, unused: true, strict: true */
/* global define, console, window, document, alert */
define(function () {

    "use strict";

    var my = {},
        conn;

    /**
     *
     * @param callback - onConnect
     * @param callback - onMessage
     */
    my.connect = function (onConnect, onMessage) {
        if (window.WebSocket) {
            // Let us open a web socket
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
                onMessage(evt);
            };
            conn.onclose = function (data) {
                console.log("connection was closed to '" + data.currentTarget.URL + '"');
            };
        } else {
            alert("Your browser does not support WebSockets. :'|");
        }
    };

    my.close = function () {
        conn.close();
    };

    /**
     *
     * @param message
     * @returns {boolean}
     */
    my.send = function (message) {
        if (conn.readyState != window.WebSocket.OPEN) {
            return false;
        }
        conn.send(message);
        return true;
    };

    // Make the function wait until the connection is made...
    return my;
});