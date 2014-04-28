define(function () {

    /**
     *
     * @type {number}
     */
    var MOVE_UP = 0,
        MOVE_DOWN = 1,
        MOVE_RIGHT = 2,
        MOVE_LEFT = 3,
        // A byte, represents the actions
        actions = 0
        // console.log the keypressed
        debug = false;

    /**
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
        37: MOVE_LEFT // arrow left
    }

    /**
     *
     * @param event
     */
    window.document.onkeydown = function (event) {
        if(debug) {
            console.log(String.fromCharCode(event.keyCode), event.keyCode);
        }
        if(typeof keycodeToAction[event.keyCode] === 'undefined') { return; }
        actions |= 1 <<  keycodeToAction[event.keyCode];
    }

    /**
     *
     * @param event
     */
    window.document.onkeyup = function (event) {
        if(typeof keycodeToAction[event.keyCode] === 'undefined') { return; }
        actions &= ~(1 << keycodeToAction[event.keyCode]);
    }

    /**
     *
     * @type {{}}
     */
    var my = {};
    my.get = function () {
        return actions;
    }
    return my;
});