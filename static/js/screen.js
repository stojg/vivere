define(function (){

    var my = {};

    function getSize() {
        // Get current screen width
        var w = window,
            d = document,
            e = d.documentElement,
            g = d.getElementsByTagName('body')[0],
            x = w.innerWidth || e.clientWidth || g.clientWidth,
            y = w.innerHeight || e.clientHeight || g.clientHeight;
        return {width: x, height:y}
    }

    my.size = getSize();

    return my;

});



