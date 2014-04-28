define(["gamestate"], function (gamestate) {

    var lastUpdated = 0;

    var sim = {}



    sim.update = function (currentTime) {

        for (var i = 0; i < gamestate.entities.length; i++) {
            if (typeof(gamestate.entities[i]) === 'undefined') {
                continue;
            }
            if (gamestate.entities[i].action == 4) {
                this.stages[0].removeChild(gamestate.entities[i]);
                delete(gamestate.entities[i]);
            }
            gamestate.entities[i].update(currentTime - lastUpdated);
            // @todo some clever lerp:ing

        }
        lastUpdated = currentTime;
    }

    return sim;
})