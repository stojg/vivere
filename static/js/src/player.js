/* jshint undef: true, unused: true, strict: true */
/* global define, console, window */
define(['src/world', 'lib/datastream'], function (world, DataStream) {

    "use strict";

    var player = {};

    var cmdSequence = 0;

    var cameraPosition = [0, 0];

    var cameraZoom = 1;

    /**
     *
     * @type {number}
     */
    var MOVE_UP = 0,
        MOVE_DOWN = 1,
        MOVE_RIGHT = 2,
        MOVE_LEFT = 3,
        ZOOM_IN = 4,
        ZOOM_OUT = 5,
    // A byte, represents the actions
        actions = 0,
    // console.log the keypressed
        debug = false;

    /**
     * http://www.javascriptkeycode.com/
     *
     * @type []int
     */
    var keycodeToAction = {
        87: MOVE_UP, // w
        38: MOVE_UP, // arrow up
        83: MOVE_DOWN, // s
        40: MOVE_DOWN, // arrow down
        68: MOVE_RIGHT, // d
        39: MOVE_RIGHT, // arrow right
        65: MOVE_LEFT, // a
        37: MOVE_LEFT, // arrow left
        81: ZOOM_OUT, // q
        69: ZOOM_IN //
    };

    /**
     *
     * @param event
     */
    window.document.onkeydown = function (event) {

        if(keycodeToAction[event.keyCode] === MOVE_LEFT){
            cameraPosition[0] += 15;
            return
        }
        if(keycodeToAction[event.keyCode] === MOVE_RIGHT){
            cameraPosition[0] -= 15;
            return
        }
        if(keycodeToAction[event.keyCode] === MOVE_UP){
            cameraPosition[1] += 15;
            return
        }
        if(keycodeToAction[event.keyCode] === MOVE_DOWN){
            cameraPosition[1] -= 15;
            return
        }

        if(keycodeToAction[event.keyCode] === ZOOM_IN){
            cameraZoom -= 0.2;
            if(cameraZoom < 0.5) {
                cameraZoom = 0.5;
            }
            return
        }

        if(keycodeToAction[event.keyCode] === ZOOM_OUT){
            cameraZoom += 0.2;
            if(cameraZoom > 10) {
                cameraZoom = 10;
            }
            return
        }

        if (debug) {
            console.log(String.fromCharCode(event.keyCode), event.keyCode);
        }
        if (typeof keycodeToAction[event.keyCode] === 'undefined') {
            return;
        }
        actions |= 1 << keycodeToAction[event.keyCode];
    };

    /**
     *
     * @param event
     */
    window.document.onkeyup = function (event) {
        if (typeof keycodeToAction[event.keyCode] === 'undefined') {
            return;
        }
        actions &= ~(1 << keycodeToAction[event.keyCode]);
    };

    /**
     * Send the client commands back to the server
     *
     * @returns bool
     */
    player.sendUpdates = function (tickLength, websocket) {

        if (actions ==+ 0) {
            return false;
        }
        var cmd = new DataStream();
        // Attach the current server time


        // send a input request
        cmd.writeUint8(3);
        // Send the latest server time recieved
        // cmd.writeFloat64(tickLength);
        // Send which # of command this is, this get passed by the server
        // so we can readjust the estimated controlled entity's position
        cmd.writeUint32(++cmdSequence);
        // How many ms this command was pressed for
        cmd.writeUint32(tickLength);
        // a byte that represents the pressed buttons
        cmd.writeUint32(actions);

        //return websocket.send(cmd.buffer);
    };

    player.camera = function() {
        return [cameraPosition, cameraZoom];
    }

    return player;
});
