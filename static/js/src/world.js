/* jshint undef: true, unused: true, strict: true */
/* global define */
define(['src/entity', 'lib/babylon.2.3.max'], function (entity) {

    "use strict";

    return function (scene) {

        var blue = new BABYLON.StandardMaterial("texture1", scene);
        blue.diffuseColor = new BABYLON.Color3(0.2, 0.2, 0.7);
        blue.specularColor = new BABYLON.Color3(0.0, 0.0, 0.0);

        var front = new BABYLON.StandardMaterial("texture1", scene);
        front.diffuseColor = new BABYLON.Color3(1.0, 0.0, 0.0);

        var pink = new BABYLON.StandardMaterial("texture1", scene);
        pink.diffuseColor = new BABYLON.Color3(1.0, 0.2, 0.7);

        var red = new BABYLON.StandardMaterial("texture1", scene);
        red.diffuseColor = new BABYLON.Color3(1.0, 0.4, 0.4);

        var moccasin = new BABYLON.StandardMaterial("texture1", scene);
        moccasin.diffuseColor = new BABYLON.Color3(1.0,0.9, 0.8);

        var pinkLight = new BABYLON.StandardMaterial("texture1", scene);
        pinkLight.diffuseColor = new BABYLON.Color3(.7, 0.3, .7);
        pinkLight.specularColor = new BABYLON.Color3(.7, 0.3, .7);

        this.templates = {};

        var box = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        box.scaling = new BABYLON.Vector3(30, 15, 30);
        box.isVisible = false;
        box.material = blue;
        this.templates[1] = box;

        var pray = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        pray.scaling = new BABYLON.Vector3(30, 30, 30);
        pray.isVisible = false;
        //pray.material = pink;
        this.templates[2] = pray;

        var yellow = new BABYLON.StandardMaterial("texture1", scene);
        yellow.diffuseColor = new BABYLON.Color3(1.0,0.9, 0.8);
        var blue = new BABYLON.StandardMaterial("texture1", scene);
        blue.diffuseColor = new BABYLON.Color3(.4, .5, 1);
        var green = new BABYLON.StandardMaterial("texture1", scene);
        green.diffuseColor = new BABYLON.Color3(.5, 1.0, .4);

        var multi=new BABYLON.MultiMaterial("nuggetman",scene);
        multi.subMaterials.push(green);
        multi.subMaterials.push(yellow);
        multi.subMaterials.push(front);
        multi.subMaterials.push(yellow);
        multi.subMaterials.push(blue);
        multi.subMaterials.push(yellow);
        pray.subMeshes=[];
        var verticesCount=pray.getTotalVertices();
        pray.subMeshes.push(new BABYLON.SubMesh(0, 0, verticesCount, 0, 6, pray));
        pray.subMeshes.push(new BABYLON.SubMesh(1, 1, verticesCount, 6, 6, pray));
        pray.subMeshes.push(new BABYLON.SubMesh(2, 2, verticesCount, 12, 6, pray));
        pray.subMeshes.push(new BABYLON.SubMesh(3, 3, verticesCount, 18, 6, pray));
        pray.subMeshes.push(new BABYLON.SubMesh(4, 4, verticesCount, 24, 6, pray));
        pray.subMeshes.push(new BABYLON.SubMesh(5, 5, verticesCount, 30, 6, pray));
        pray.material=multi;


        var hunter = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        hunter.scaling = new BABYLON.Vector3(30, 30, 30);
        hunter.isVisible = false;
        hunter.material = red;
        this.templates[3] = hunter;

        var scared = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        scared.scaling = new BABYLON.Vector3(10, 10, 10);
        scared.isVisible = false;
        scared.material = moccasin;
        this.templates[4] = scared;

        var taken = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        taken.scaling = new BABYLON.Vector3(10, 10, 10);
        taken.isVisible = false;
        taken.material = pinkLight;
        this.templates[5] = taken;

        this.entities = [];

        this.serverTick = 0;

        this.update = function (buf, scene) {
            // first byte is current servertick
            this.serverTick = buf.readFloat32();

            var commands = [];

            var id = 0;
            while (!buf.isEof()) {
                var cmd = buf.readUint8();
                switch (cmd) {
                    // INST_ENTITY_ID
                    case 1:
                        // we are changing entity, set a new ID that will be used by all the following non 1 cmds
                        id = buf.readFloat32();
                        if (typeof this.entities[id] == 'undefined') {
                            this.entities[id] = entity.create(id, 120, scene, this.templates);
                        }
                        commands[id] = {};
                        commands[id].timestamp = window.performance.now();
                        break;
                    // INST_SET_POSITION
                    case 2:
                        commands[id].position = {x: buf.readFloat32(), y: buf.readFloat32(), z: buf.readFloat32()};
                        break;
                    // INST_SET_ROTATION
                    case 3:
                        commands[id].orientation = [];
                        commands[id].orientation[0] = buf.readFloat32();
                        commands[id].orientation[1] = buf.readFloat32();
                        commands[id].orientation[2] = buf.readFloat32();
                        commands[id].orientation[3] = buf.readFloat32();
                        break;
                    // INST_SET_MODEL
                    case 4:
                        commands[id].model = buf.readFloat32();
                        break;
                    // INST_SET_SCALE
                    case 5:
                        commands[id].scale = {x: buf.readFloat32(), y: buf.readFloat32(), z: buf.readFloat32()};
                        break;

                }
            }

            for (id in commands) {
                this.entities[id].serverUpdate(commands[id]);
            }
        }
    };
});
