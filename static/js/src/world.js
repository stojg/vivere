/* jshint undef: true, unused: true, strict: true */
/* global define */
define(['src/entity', 'lib/babylon.2.3.max'], function (entity) {

    "use strict";

    return function (scene) {

        var blue_material = new BABYLON.StandardMaterial("texture1", scene);
        blue_material.diffuseColor = new BABYLON.Color3(0.0, 0.0, 0.4);

        var red_material = new BABYLON.StandardMaterial("red_material", scene);
        red_material.diffuseColor = new BABYLON.Color3(0.9, 0.2, 0.2);

        var pink_material = new BABYLON.StandardMaterial("texture1", scene);
        pink_material.diffuseColor = new BABYLON.Color3(1.0, 0.2, 0.7);

        var moccasin = new BABYLON.StandardMaterial("texture1", scene);
        moccasin.diffuseColor = new BABYLON.Color3(1.0,0.9, 0.8);

        var pinkLight = new BABYLON.StandardMaterial("texture1", scene);
        pinkLight.diffuseColor = new BABYLON.Color3(.7, 0.3, .7);

        var yellow_material = new BABYLON.StandardMaterial("yellow", scene);
        yellow_material.diffuseColor = new BABYLON.Color3(0.9, 0.8, 0.7);

        var blue_material = new BABYLON.StandardMaterial("texture1", scene);
        blue_material.diffuseColor = new BABYLON.Color3(.5, .5, 0.9);

        var green_material = new BABYLON.StandardMaterial("texture1", scene);
        green_material.diffuseColor = new BABYLON.Color3(.5, 1.0, .4);

        var ground_material = new BABYLON.StandardMaterial("texture1", scene);
        ground_material.diffuseColor = new BABYLON.Color3(0.2, 0.2, 0.3);
        ground_material.specularColor = new BABYLON.Color3(0.0, 0.0, 0.0);

        this.templates = {};

        var ground = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        ground.scaling = new BABYLON.Vector3(1, 1, 1);
        ground.isVisible = false;
        ground.material = ground_material;
        this.templates[1] = ground;

        var box_model = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        box_model.scaling = new BABYLON.Vector3(30, 15, 30);
        box_model.isVisible = false;
        box_model.material = blue_material;
        this.templates[2] = box_model;

        var pray_model = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        pray_model.scaling = new BABYLON.Vector3(30, 30, 30);
        pray_model.isVisible = false;
        this.templates[3] = pray_model;

        var multi=new BABYLON.MultiMaterial("nuggetman",scene);
        //multi.subMaterials.push(green);
        multi.subMaterials.push(yellow_material);
        multi.subMaterials.push(yellow_material);
        multi.subMaterials.push(red_material);
        multi.subMaterials.push(yellow_material);
        multi.subMaterials.push(yellow_material);
        multi.subMaterials.push(yellow_material);
        pray_model.subMeshes=[];
        var verticesCount=pray_model.getTotalVertices();
        pray_model.subMeshes.push(new BABYLON.SubMesh(0, 0, verticesCount, 0, 6, pray_model));
        pray_model.subMeshes.push(new BABYLON.SubMesh(1, 1, verticesCount, 6, 6, pray_model));
        pray_model.subMeshes.push(new BABYLON.SubMesh(2, 2, verticesCount, 12, 6, pray_model));
        pray_model.subMeshes.push(new BABYLON.SubMesh(3, 3, verticesCount, 18, 6, pray_model));
        pray_model.subMeshes.push(new BABYLON.SubMesh(4, 4, verticesCount, 24, 6, pray_model));
        pray_model.subMeshes.push(new BABYLON.SubMesh(5, 5, verticesCount, 30, 6, pray_model));
        pray_model.material=multi;

        var hunter = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        hunter.scaling = new BABYLON.Vector3(30, 30, 30);
        hunter.isVisible = false;
        hunter.material = red_material;
        this.templates[4] = hunter;

        var scared = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        scared.scaling = new BABYLON.Vector3(10, 10, 10);
        scared.isVisible = false;
        scared.material = moccasin;
        this.templates[5] = scared;

        var taken = BABYLON.Mesh.CreateBox("box", 1.0, scene, false, BABYLON.Mesh.DEFAULTSIDE);
        //var sphere = BABYLON.Mesh.CreateSphere("sphere", 20, 1.0, scene);
        taken.scaling = new BABYLON.Vector3(10, 10, 10);
        taken.isVisible = false;
        taken.material = pinkLight;
        this.templates[6] = taken;

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
                            this.entities[id] = entity.create(id, 100, scene, this.templates);
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
