require(["screen", "websocket", 'pixi', 'entity', "gamestate", "commands", "simulator"], function (screen, websocket, pixi, entity, gamestate, commands, simulator) {

    var main = {};

    // Simulation tick time in (1000 / 20hz = 50ms)
    main.tickLength = 50;
    main.cmdSequence = 0;
    main.connected = false;
    main.stopGameLoop = 0;
    main.lastTick = window.performance.now();
    main.lastRender = main.lastTick;
    main.pixi = null;
    main.commandTick = 0;
    main.stages = [];

    /**
     * Initialize the renderer and the gamestate
     */
    main.init = function () {
        this.pixi = pixi.autoDetectRenderer(1000, 600);
        document.body.appendChild(this.pixi.view);
        this.stages[0] = new pixi.Stage(0x666666);
    }

    /**
     * Render the game
     */
    main.render = function () {
        for(var i = 0; i < this.stages.length; i++) {
           this.pixi.render(this.stages[i]);
        }
    }

    /**
     * Send the client commands back to the server
     *
     * @returns bool
     */
    main.sendUpdates = function(tickLength) {
        if(commands.get() == 0) {
            return false;
        }
        var cmd = new DataStream();
        cmd.writeUint32(gamestate.serverTick);
        cmd.writeUint32(++main.cmdSequence);
        cmd.writeUint32(tickLength);
        cmd.writeUint32(commands.get());
        return websocket.send(cmd.buffer);
    }

    /**
     * Behold, the game server starts after the websocket connects
     */
    websocket.connect(function () {
        main.connected = true;
        main.init();
        gameloop();
    }, onRecieve);

    /**
     * The main game loop
     *
     * @param tFrame - high resolution timer
     */
    function gameloop(tFrame) {
        main.stopGameLoop = window.requestAnimationFrame(gameloop);
        var nextTick = main.lastTick + main.tickLength;
        var numTicks = 0;
        if (tFrame > nextTick) {
            var timeSinceTick = tFrame - main.lastTick;
            numTicks = Math.floor(timeSinceTick / main.tickLength);
        }
        for (var i = 0; i < numTicks; i++) {
            main.lastTick = main.lastTick + main.tickLength;
            simulator.update(main.lastTick);
        }
        main.render();
        main.sendUpdates(main.tickLength);
        main.lastRender = tFrame;
    }

    /**
     * Gets called by the websocket when things are happening
     *
     * Updates the gamestate with the data from the server
     *
     * @param evt
     */
    function onRecieve(evt) {
        var buf = new DataStream(evt.data)
        gamestate.serverTick = buf.readUint32();
        var nEnts = buf.readUint16();
        for (i = 0; i < nEnts; i++) {
            // get the bitmask
            var bitMask = buf.readUint8();
            // id
            var id = buf.readUint16();
            // model
            if ((bitMask & (1 << 0)) > 0) {
                var modelId = buf.readUint16();
                if (modelId == 0) {
                    if (typeof gamestate.entities[id] !== 'undefined') {
                        main.stages[0].removeChild(gamestate.entities[id].getSprite());
                        delete gamestate.entities[id];
                    }
                } else if (typeof gamestate.entities[id] === 'undefined') {
                    gamestate.entities[id] = entity.create(modelId);
                    main.stages[0].addChild(gamestate.entities[id].getSprite());
                }
            }

            var command = { id: id };

            // rotation
            if ((bitMask & (1 << 1)) > 0) {
                command.rotation = buf.readFloat32();
            }
            // angular velocity
            if ((bitMask & (1 << 2)) > 0) {
                command.angularVel = buf.readFloat32();
            }
            // pos
            if ((bitMask & (1 << 3)) > 0) {
                var pos = buf.readFloat64Array(2);
                command.position = {x: pos[0], y: pos[1]};
            }
            // vel
            if ((bitMask & (1 << 4)) > 0) {
                var vel = buf.readFloat64Array(2);
                command.velocity = {x: vel[0], y: vel[1]};
            }
            // size
            if ((bitMask & (1 << 5)) > 0) {
                var size = buf.readFloat64Array(2);
                command.size = {x: size[0], y: size[1]};
            }
            // action
            if ((bitMask & (1 << 6)) > 0) {
                command.action = buf.readUint16();
            }

            gamestate.entities[id].serverUpdate(command);
        }
    }
});