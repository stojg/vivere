define(["gamestate"], function (gamestate) {

    var sim = {}

    /**
     * this method is called every approx 16ms
     *
     * @param mSec
     */
    sim.update = function(tFrame) {
        for (var i = 0; i < gamestate.entities.length; i++) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            // died
            if (gamestate.entities[i].action == 4) {
                this.stages[0].removeChild(gamestate.entities[i]);
                delete(gamestate.entities[i]);
            } else {
                gamestate.entities[i].update(tFrame);
            }
        }
    }
    return sim;
})