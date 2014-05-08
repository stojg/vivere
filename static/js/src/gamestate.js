/* jshint undef: true, unused: true, strict: true */
/* global define */
define(['src/entity'], function (entity) {

    "use strict";

    var my = {};

    my.entities = [];

    my.messageQueue = [];

    my.serverTick = 0;


    my.update = function(buf, main) {

        // first byte is current servertick
        my.serverTick = buf.readFloat32();

        var literal;

        var id = 0;
        var command = {};
        command.timestamp = 0;

        while(!buf.isEof()) {
            switch(literal = buf.readUint8()) {
                // INST_ENTITY_ID
                case 1:
                    // we are changing entity, send update to previous entity
                    if(command.timestamp !== 0) {
                        my.entities[id].serverUpdate(command);
                        command = {};
                    }
                    id = buf.readFloat32();
                    if(typeof my.entities[id] == 'undefined') {
                        my.entities[id] = entity.create(2, 120 );
                        main.stages[0].addChild(my.entities[id].getSprite());
                    }
                    command.timestamp = window.performance.now();
                    break;
                // INST_SET_POSITION
                case 2:
                    command.position = {x: buf.readFloat32(), y: buf.readFloat32()};
                    break;
                // INST_SET_ROTATION
                case 3:
                    command.rotation = buf.readFloat32();
                    break;
            }
        }
    };

    return my;


});
