/* jshint undef: true, unused: true, strict: true */
/* global require, window, clearTimeout, document */
require(["src/server", 'src/entity', "src/world", "src/player", 'lib/datastream', 'lib/babylon.2.3.max'], function (server, entity, world, player, DataStream) {

    "use strict";

    var isActive = true;
    window.onfocus = function () {
        isActive = true;
    };

    window.onblur = function () {
        isActive = false;
    };

    // Get the canvas element from our HTML above
    var canvas = document.getElementById("renderCanvas");

    // Load the BABYLON 3D engine
    var engine = new BABYLON.Engine(canvas, true);

    var main = {};

    /**
     * Initialize the renderer and the gamestate
     */
    main.init = function () {
        main.scene = this.createScene();
        main.world = new world(main.scene);
        window.addEventListener("resize", function () {
            engine.resize();
        });
    };

    main.world = {};

    main.createScene = function () {

        var scene = new BABYLON.Scene(engine);
        scene.clearColor = new BABYLON.Color3(0.1, 0.1, 0.13);
        scene.ambientColor = new BABYLON.Color3(0.3, 0.3, 0.3);
        //scene.fogMode = BABYLON.Scene.FOGMODE_LINEAR;
        //scene.fogStart = 500.0;
        //scene.fogEnd = 4000.0;
        //scene.fogColor = new BABYLON.Color3(0.0, 0.0, 0.0);
        scene.collisionsEnabled = true;

        scene.activeCamera = new BABYLON.FreeCamera("FreeCamera", new BABYLON.Vector3(1, 1, 1), scene);
        scene.activeCamera.attachControl(canvas);
        scene.activeCamera.keysUp.push(87);
        scene.activeCamera.keysLeft.push(65);
        scene.activeCamera.keysDown.push(83);
        scene.activeCamera.keysRight.push(68);
        scene.activeCamera.speed = 40;
        scene.activeCamera.checkCollisions = true;
        scene.activeCamera.setTarget(BABYLON.Vector3.Zero());
        scene.activeCamera.position = new BABYLON.Vector3(1661, 1050, 1500);
        scene.activeCamera.attachControl(canvas, false);

        //scene.debugLayer.show();

        var beforeRenderFunction = function () {
            // Camera
            if (scene.activeCamera.position.y < 30)
                scene.activeCamera.position.y = 30;
        };

        scene.registerBeforeRender(beforeRenderFunction);

        // This creates a light, aiming 0,1,0 - to the sky (non-mesh)
        //var light = new BABYLON.HemisphericLight("light1", new BABYLON.Vector3(0, 1, 0), scene);
        // Default intensity is 1. Let's dim the light a small amount
        //var light0 = new BABYLON.PointLight("Omni0", new BABYLON.Vector3(200, 100, -100), scene);
        //light0.diffuse = new BABYLON.Color3(.5, .5, .5);
        //light0.specular = new BABYLON.Color3(0.0, 0.0, 0.0);
        //light0.intensity = 1;

        var light1 = new BABYLON.HemisphericLight("Hemi0", new BABYLON.Vector3(0, 1, 0), scene);
        light1.diffuse = new BABYLON.Color3(1,1,1);
        light1.specular = new BABYLON.Color3(1, 1, 1);
        light1.groundColor = new BABYLON.Color3(0, 0, 0);
        light1.intensity = 0.95;

        //var box = BABYLON.Mesh.CreateBox("box", 100, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var shadowGenerator = new BABYLON.ShadowGenerator(1024, light0);
        //shadowGenerator.getShadowMap().renderList.push(box);

        var ground = BABYLON.Mesh.CreateGround("ground1", 3232, 3232, 4, scene);
        var groundMat = new BABYLON.StandardMaterial("texture1", scene);
        groundMat.diffuseColor = new BABYLON.Color3(0.2, 0.21, 0.21);
        //groundMat.ambientColor = new BABYLON.Color3(0.3, 0.4, 0.3);
        ground.material = groundMat;
        ground.receiveShadows = true;
        ground.checkCollisions = true;


        // Leave this function
        return scene;

    };  // End of createScene function

    main.connected = false;
    main.lastTick = window.performance.now();
    main.pixi = null;
    main.commandTick = 0;
    main.stage = [];

    main.scene = null;



    /**
     * Behold, the game server starts after the websocket connects
     */
    server.connect(function () {
        main.connected = true;
        main.init();
        engine.runRenderLoop(function () {
            var tFrame = window.performance.now();

            for (var i in main.world.entities) {
                if (typeof(main.world.entities[i]) === 'undefined') {
                    continue;
                }
                if (main.world.entities[i].state == 1) {
                    //main.stage.removeChild(world.entities[i].getSprite());
                    //world.entities.splice(i, 1)
                } else {
                    main.world.entities[i].update(tFrame);
                }
            }
            main.scene.render();
            main.lastTick = tFrame;
        });
    }, onRecieve);

    /**
     * Gets called by the websocket when things
     *
     * @param buf
     */
    function onRecieve(buf) {
        if(!isActive) {
            return;
        }
        var msgType = buf.readUint8();
        // world state update
        if (msgType === 1) {
            main.world.update(buf, main.scene)
        }
        // respond to a ping request
        if (msgType === 2) {
            server.send(server.newMessage(2));
        }
    }
});