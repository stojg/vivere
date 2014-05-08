/* jshint undef: true, unused: true, strict: true */
/* global require, window, clearTimeout, document */
require(["src/websocket", 'lib/pixi', 'src/entity', "src/gamestate", "src/player", "src/simulator", 'lib/datastream'], function (websocket, pixi, entity, gamestate, player, simulator, DataStream) {

    "use strict";

    window.cancelRequestAnimFrame = (function () {
        return window.cancelAnimationFrame ||
            window.webkitCancelRequestAnimationFrame ||
            window.mozCancelRequestAnimationFrame ||
            window.oCancelRequestAnimationFrame ||
            window.msCancelRequestAnimationFrame ||
            clearTimeout;
    })();

    var main = {};

    main.connected = false;
    main.stopGameLoop = 0;
    main.lastTick = window.performance.now();
    main.lastRender = window.performance.now();
    main.pixi = null;
    main.commandTick = 0;
    main.stages = [];

    main.fpsText = null;
    main.frameCounter = 0;

    main.mpsText = null;
    main.messageCounter = 0;
    main.lastRecieved = window.performance.now();

    /**
     * Initialize the renderer and the gamestate
     */
    main.init = function () {
        this.pixi = pixi.autoDetectRenderer(1000, 600);
        document.body.appendChild(this.pixi.view);
        this.stages[0] = new pixi.Stage(0x666666);
        main.fpsText = new pixi.Text("fps ", {font: "22px Arial", fill: "white"});
        main.fpsText.position = {x: 10, y: 5};
        this.stages[0].addChild(main.fpsText);
        main.mpsText = new pixi.Text("mps ", {font: "22px Arial", fill: "white"});
        main.mpsText.position = {x: 10, y: 25};
        this.stages[0].addChild(main.mpsText);
        this.lastTick = window.performance.now();

//        setTimeout(function(){
//            window.cancelRequestAnimFrame(main.stopGameLoop);
//            websocket.close();
//        }, 1.1*1000)

    };

    /**
     * Render the game
     */
    main.render = function () {
        main.frameCounter++;
        for (var i = 0; i < this.stages.length; i++) {
            this.pixi.render(this.stages[i]);
        }
    };

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
        var elapsed = tFrame - main.lastTick;
        simulator.update(tFrame, main);
        main.render();
        updateFPSCounter(tFrame);
        updateMPSCounter();
        player.sendUpdates(elapsed, websocket);
        main.lastTick = tFrame;
    }

    /**
     * Prints frames rendered per second
     *
     * @param t
     */
    function updateFPSCounter(tFrame) {
        main.fpsText.setText("fps " + Math.round(1000 / (tFrame / main.frameCounter)));
    }

    /**
     * Prints messages recieved per second
     */
    function updateMPSCounter() {
        main.mps = 1000 / (main.lastRecieved / main.messageCounter);
        main.mpsText.setText("mps " + Math.round(main.mps));
    }

    /**
     * Gets called by the websocket when things
     *
     * @param evt
     */
    function onRecieve(buf) {

        var msgType = buf.readUint8();
        // world state update
        if(msgType === 1) {
            gamestate.update(buf, main)
        }
        // respond to a ping request
        if(msgType === 2) {
            websocket.send(websocket.newMessage(2));
        }

        return;
    }
});