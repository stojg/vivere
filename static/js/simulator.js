define(["gamestate"], function (gamestate) {

    var sim = {}

    /**
     * This method is run every 50ms
     *
     * @param currentTime
     */
    sim.applyUpdates = function (currentTime) {
        for (var i = 0; i < gamestate.entities.length; i++) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            if (gamestate.entities[i].action == 4) {
                this.stages[0].removeChild(gamestate.entities[i]);
                delete(gamestate.entities[i]);
            }
            gamestate.entities[i].applyUpdates();
        }
    }

    /**
     * this method is called every approx 16ms
     *
     * @param mSec
     */
    sim.update = function(mSec) {
        for (var i = 0; i < gamestate.entities.length; i++) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            gamestate.entities[i].update(mSec);
        }
    }

    /**
     * this method is called every approx 16ms
     *
     * @param msec
     */
    sim.interpolate = function(range) {
        for (var i = 0; i < gamestate.entities.length; i++) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            gamestate.entities[i].interpolate(range);
        }
    }

    return sim;
})