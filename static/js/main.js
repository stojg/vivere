require(["screen", "websocket", 'pixi', 'entity', "gamestate", "commands", "simulator"], function (screen, websocket, pixi, entity, gamestate, commands, simulator) {

    window.cancelRequestAnimFrame = ( function() {
        return window.cancelAnimationFrame          ||
            window.webkitCancelRequestAnimationFrame    ||
            window.mozCancelRequestAnimationFrame       ||
            window.oCancelRequestAnimationFrame     ||
            window.msCancelRequestAnimationFrame        ||
            clearTimeout
    } )();

    var main = {};

    // Simulation tick time in (1000 / 60hz = 50ms)
    main.tickLength = 17;
    main.cmdSequence = 0;
    main.connected = false;
    main.stopGameLoop = 0;
    main.lastTick = window.performance.now();
    main.lastRender = window.performance.now();
    main.pixi = null;
    main.commandTick = 0;
    main.stages = [];

    main.fpsText;
    main.frameCounter = 0;

    main.mpsText;
    main.messageCounter = 0;
    main.lastRecieved = window.performance.now();

    /**
     * Initialize the renderer and the gamestate
     */
    main.init = function () {
        this.pixi = pixi.autoDetectRenderer(1000, 600);
        document.body.appendChild(this.pixi.view);
        this.stages[0] = new pixi.Stage(0x666666);
        main.fpsText = new pixi.Text("fps ", {font:"22px Arial", fill:"white"});
        this.stages[0].addChild(main.fpsText);
        main.mpsText = new pixi.Text("mps ", {font:"22px Arial", fill:"white"});
        main.mpsText.position = {x:0, y:20}
        this.stages[0].addChild(main.mpsText);
        this.lastTick = window.performance.now();
        this.lastRender = window.performance.now();

//        setTimeout(function(){
//            window.cancelRequestAnimFrame(main.stopGameLoop);
//            websocket.close();
//        }, 5*1000)

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
        gameloop(window.performance.now());
    }, onRecieve);

    main.lastRenderTime = 0;

    /**
     * The main game loop
     *
     * @param tFrame - high resolution timer
     */
    function gameloop(tFrame) {
        main.stopGameLoop = window.requestAnimationFrame(gameloop);
        simulator.update(tFrame);
        main.render();
        printFPS(tFrame);
        main.lastRender = tFrame;
        main.sendUpdates(tFrame);
    }

    /**
     * Prints frames rendered per second
     *
     * @param t
     */
    function printFPS(tFrame) {
        main.frameCounter++;
        main.fpsText.setText("fps " + Math.round(1000 / (tFrame / main.frameCounter)));
    }

    /**
     * Prints messages recieved per second
     */
    function printMPS() {
        main.messageCounter++;
        main.mps = 1000 / (main.lastRecieved / main.messageCounter);
        main.mpsText.setText("mps " + Math.round(main.mps));
    }

    /**
     * Gets called by the websocket when things
     *
     * @param evt
     */
    function onRecieve(evt) {
        if(gamestate.firstSnapshotTS === 0) {
            gamestate.firstSnapshotTS = window.performance.now();
        }

        main.lastRecieved = window.performance.now();
        printMPS();

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
            // @todo subtract network lag (RTT) to this
            command.timestamp = window.performance.now();

            command.tick = gamestate.serverTick;

            // rotation
            if ((bitMask & (1 << 1)) > 0) {
                command.rotation = buf.readFloat64();
            }

            // pos
            if ((bitMask & (1 << 2)) > 0) {
                var pos = buf.readFloat64Array(2);
                command.position = {x: pos[0], y: pos[1]};
            }

            // vel
            if ((bitMask & (1 << 3)) > 0) {
                var vel = buf.readFloat64Array(2);
                command.velocity = {x: vel[0], y: vel[1]};
            }
            // size
            if ((bitMask & (1 << 4)) > 0) {
                var size = buf.readFloat64Array(2);
                command.size = {x: size[0], y: size[1]};
            }
            // action
            if ((bitMask & (1 << 5)) > 0) {
                command.action = buf.readUint16();
            }

            gamestate.entities[id].serverUpdate(command);
        }
    }
});