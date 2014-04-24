define(function () {

    // Find key codes here: http://unixpapa.com/js/testkey.html

    // A byte representation of binaries
    var actions = 0;

    // Add commands to the actions, bit combined
    window.document.onkeydown = function(event) {
        var key_press = String.fromCharCode(event.keyCode);
        var key_code = event.keyCode;
        // Up
        if (key_code == 87 || key_code == 38) {
            actions |= 1<<0
        }
        // Down
        if (key_code == 83 || key_code == 40) {
            actions |= 1<<1
        }
        // Right
        if (key_code == 68 || key_code == 39) {
            actions |= 1<<2
        }
        // Left
        if (key_code == 65 || key_code == 37) {
            actions |= 1<<3
        }
    }
    window.document.onkeyup = function(event){
        var key_press = String.fromCharCode(event.keyCode);
        var key_code = event.keyCode;
        // unset bit one
        if (key_code == 87 || key_code == 38) {
            actions &= ~(1<<0)
        }
        // unset bit two
        if (key_code == 83 || key_code == 40) {
            actions &= ~(1<<1)
        }
        // set bit three
        if (key_code == 68 || key_code == 39) {
            actions &= ~(1<<2)
        }
        // set bit four
        if (key_code == 65 || key_code == 37) {
            actions &= ~(1<<3)
        }
    }

    var my = {};

    my.get = function() {
        return actions;
    }

    return my;
});