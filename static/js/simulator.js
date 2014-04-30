define(["gamestate"], function (gamestate) {

    var sim = {}

    /**
     * this method is called every approx 16ms
     *
     * @param mSec
     */
    sim.update = function(tFrame, main) {
        for(i in gamestate.entities) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            if (gamestate.entities[i].state == 1) {
                main.stages[0].removeChild(gamestate.entities[i].getSprite());
                delete gamestate.entities[i];
            } else {
                gamestate.entities[i].update(tFrame);
            }
        }
    }
    return sim;
})