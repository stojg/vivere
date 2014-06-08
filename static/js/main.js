/* jshint undef: true, unused: true, strict: true */
/* global require, window, clearTimeout, document */
require(["src/server", 'lib/pixi', 'src/entity', "src/world", "src/player", "src/simulator", 'lib/datastream'], function (server, pixi, entity, world, player, simulator, DataStream) {

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
    main.lastTick = window.performance.now();
    main.pixi = null;
    main.commandTick = 0;
    main.stage = [];

    /**
     * Initialize the renderer and the gamestate
     */
    main.init = function () {
        this.pixi = pixi.autoDetectRenderer(1024, 640);
        document.body.appendChild(this.pixi.view);
        this.stage = new pixi.Stage(0x666666);
        this.stage.addChild(world.container);
        this.lastTick = window.performance.now();
//        setTimeout(function(){
//            window.cancelRequestAnimFrame(main.stopGameLoop);
//            websocket.close();
//        }, 1.1*1000)

    };

    /**
     * Render the game
     */
    main.render = function (elapsed) {
        var camera = player.camera();
        main.frameCounter++;
        // https://github.com/anvaka/ngraph/blob/master/examples/pixi.js/03%20-%20Zoom%20And%20Pan/globalInput.js
        world.container.position.x = camera[0][0];
        world.container.position.y = camera[0][1];
        world.container.scale.x = 1 / camera[1];
        world.container.scale.y = 1 / camera[1];
        this.pixi.render(this.stage);
    };

    /**
     * Behold, the game server starts after the websocket connects
     */
    server.connect(function () {
        main.connected = true;
        main.init();
        gameloop(window.performance.now());
    }, onRecieve);

    /**
     * The main game loop
     *
     * @param tFrame - high resolution timer
     */
    function gameloop(tFrame) {
        main.stopGameLoop = window.requestAnimationFrame(gameloop);
        var elapsed = tFrame - main.lastTick;
        simulator.update(tFrame, main);
        main.render(tFrame);
        player.sendUpdates(elapsed, server);
        main.lastTick = tFrame;
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
            world.update(buf, main)
        }
        // respond to a ping request
        if(msgType === 2) {
            server.send(server.newMessage(2));
        }

        return;
    }
});