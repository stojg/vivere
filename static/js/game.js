"use strict";

require(["src/server", 'src/entity', "src/world"], function (server, entity, world) {

    var info = document.getElementById("info");
    info.innerHTML = "initializing";

    var isActive = true;
    window.onfocus = function () {
        info.innerHTML = "";
        isActive = true;
    };

    window.onblur = function () {
        info.innerHTML = "paused";
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
        info.innerHTML = "";
        window.addEventListener("resize", function () {
            engine.resize();
        });
    };

    main.world = {};


    main.createScene = function () {

        var scene = new BABYLON.Scene(engine);
        scene.clearColor = new BABYLON.Color3(0.1, 0.1, 0.13);
        scene.ambientColor = new BABYLON.Color3(0.3, 0.3, 0.3);
        scene.collisionsEnabled = true;

        scene.activeCamera = new BABYLON.FreeCamera("FreeCamera", new BABYLON.Vector3(1, 1, 1), scene);
        scene.activeCamera.attachControl(canvas);
        scene.activeCamera.keysUp.push(87);
        scene.activeCamera.keysLeft.push(65);
        scene.activeCamera.keysDown.push(83);
        scene.activeCamera.keysRight.push(68);
        scene.activeCamera.speed = 40;
        scene.activeCamera.checkCollisions = true;
        scene.activeCamera.position = new BABYLON.Vector3(500, 200, -500);
        scene.activeCamera.setTarget(new BABYLON.Vector3(-500, -200,  500));
        scene.activeCamera.attachControl(canvas, false);

        // debug axis indicators
        //var originBox = BABYLON.Mesh.CreateBox("box", 1.0, scene);
        //originBox.scaling = new BABYLON.Vector3(30, 5, 30);
        //originBox.position = new BABYLON.Vector3(0, -2, 0);
        //var boxXpos = BABYLON.Mesh.CreateBox("box", 1.0, scene);
        //boxXpos.scaling = new BABYLON.Vector3(10, 10, 10);
        //boxXpos.position.x = 300;
        //boxXpos.position.y = 0;
        //boxXpos.position.z = 0;
        //var red = new BABYLON.StandardMaterial("texture1", scene);
        //red.diffuseColor = new BABYLON.Color3(1.0, .2, .2);
        //boxXpos.material = red;
        //var boxYpos = boxXpos.clone();
        //boxYpos.position = {x: 0, y: 100, z: 0};
        //var blue = new BABYLON.StandardMaterial("texture1", scene);
        //blue.diffuseColor = new BABYLON.Color3(.4, .5, 1);
        //boxYpos.material = blue ;
        //var boxZpos = boxXpos.clone();
        //boxZpos.position = {x: 0, y: 0, z: 300};
        //var green = new BABYLON.StandardMaterial("texture1", scene);
        //green.diffuseColor = new BABYLON.Color3(.5, 1.0, .4);
        //boxZpos.material = green;

        //scene.debugLayer.show();

        var beforeRenderFunction = function () {
            // Camera
            if (scene.activeCamera.position.y < 30)
                scene.activeCamera.position.y = 30;
        };

        scene.registerBeforeRender(beforeRenderFunction);

        var light1 = new BABYLON.HemisphericLight("Hemi0", new BABYLON.Vector3(0.3, 1, -1), scene);
        light1.diffuse = new BABYLON.Color3(1,1,1);
        light1.specular = new BABYLON.Color3(1, 1, 1);
        light1.groundColor = new BABYLON.Color3(0, 0, 0);
        light1.intensity = 0.95;

        return scene;

    };  // End of createScene function

    main.connected = false;
    main.lastTick = window.performance.now();
    main.pixi = null;
    main.commandTick = 0;
    main.stage = [];

    main.scene = null;

    info.innerHTML = "connecting";
    /**
     * Behold, the game server starts after the websocket connects
     */
    server.connect(function () {
        main.connected = true;
        info.innerHTML = "init";
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
