/* jshint undef: true, unused: true, strict: true */
/* global define */
define(["src/world"], function (world) {

    "use strict";

    var sim = {};

    /**
     * this method is called every approx 16ms
     *
     * @param mSec
     */
    sim.update = function (tFrame, main) {
        for (var i in world.entities) {

            if (typeof(world.entities[i]) === 'undefined') {
                continue;
            }

            if (world.entities[i].state == 1) {
                main.stages[0].removeChild(world.entities[i].getSprite());
                world.entities.splice(i, 1)
            } else {
                world.entities[i].update(tFrame);
            }
        }
    };

    return sim;
});