/* jshint undef: true, unused: true, strict: true */
/* global define */
define(['src/entity', "lib/pixi"], function (entity, PIXI) {

    "use strict";

    var world = {};

    world.entities = [];

    world.container = new PIXI.DisplayObjectContainer();

    world.container.pivot = {x: -1024/2, y: -640/2};

    world.camera = new PIXI.DisplayObjectContainer();

    world.serverTick = 0;

    world.update = function(buf, main) {
        // first byte is current servertick
        world.serverTick = buf.readFloat32();

        var id = 0;

        var commands = [];

        while(!buf.isEof()) {
            switch(buf.readUint8()) {
                // INST_ENTITY_ID
                case 1:
                    // we are changing entity, send update to previous entity
                    id = buf.readFloat32();
                    if(typeof world.entities[id] == 'undefined') {
                        world.entities[id] = entity.create(2, 120);
                        this.container.addChild(world.entities[id].getSprite());
                    }
                    commands[id] = {}
                    commands[id].timestamp = window.performance.now();
                    break;
                // INST_SET_POSITION
                case 2:
                    commands[id].position = {x: buf.readFloat32(), y: buf.readFloat32()};
                    break;
                // INST_SET_ROTATION
                case 3:
                    commands[id].orientation = buf.readFloat32();
                    break;
                // INST_SET_MODEL
                case 4:
                    commands[id].model = buf.readFloat32();
                    break;
            }
        }
        for (id in commands) {
            world.entities[id].serverUpdate(commands[id]);
        }
    };
    return world;
});
