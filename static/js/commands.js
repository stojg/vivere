define(function () {

    // Find key codes here: http://unixpapa.com/js/testkey.html

    // A byte representation of binaries
    var actions = 0;

    // Add commands to the actions, bit combined
    window.document.onkeydown = function(event) {
        var key_press = String.fromCharCode(event.keyCode);
        var key_code = event.keyCode;
        // set bit one
        if (key_code == 87) {
            actions |= 1<<0
        }
        // set bit one
        if (key_code == 83) {
            actions |= 1<<1
        }
    }
    window.document.onkeyup = function(event){
        var key_press = String.fromCharCode(event.keyCode);
        var key_code = event.keyCode;
        // unset bit one
        if (key_code == 87) {
            actions &= ~(1<<0)
        }
        // unset bit two
        if (key_code == 83) {
            actions &= ~(1<<1)
        }
    }

    var my = {};

    my.get = function() {
        return actions;
    }

    return my;
});