/* jshint undef: true, unused: true, strict: true */
/* global define, console */
define(["lib/babylon.2.3.max"], function () {

    'use strict';

    var GameObject = function (id, scene, templates) {

        this.id = id;

        this.scene = scene;

        this.model = null;

        this.interpolationDelay = 0;

        this.sprite = null;

        this.templates = templates;

        /**
         *
         * @type {number}
         */
        this.state = 0;

        /**
         * Contains a list of queued updates from the server
         *
         * Older update -> last update
         *
         * @type {Array}
         */
        this.server = [];

        /**
         * Contains a list of updates in the past, used for interpolation
         * Last snapshot -> older
         *
         * @type {Array}
         */
        this.snapshots = [];

        /**
         * serverUpdate is called when the server sends an update command
         * to this entity
         *
         * @param message
         */
        this.serverUpdate = function (message) {
            this.server.push(message);
            // biggest size of the queue is 20 history items, roughtly one sec
            if (this.server.length > 20) {
                this.server.unshift();
            }
        };

        /**
         *
         */
        this.applyServerUpdates = function () {
            // Move queued server updates to the snapshot array
            var msg = this.server.pop();
            while (typeof msg !== 'undefined') {
                this.snapshots.unshift(msg);
                msg = this.server.pop();
            }
        };

        /**
         *
         * @param tFrame
         */
        this.update = function (tFrame) {
            var coef,
                latestSnapshot,
                interpolationTime = tFrame - this.interpolationDelay;

            if (typeof tFrame == 'undefined' || tFrame === 0) {
                return;
            }


            this.applyServerUpdates();

            latestSnapshot = this.getLatestState(tFrame);
            if (latestSnapshot === false) {
                return false;
            }

            this.state = latestSnapshot.state;
            //this.sprite.rotation = latestSnapshot.orientation;

            if (this.model != latestSnapshot.model) {
                if(this.sprite) {
                    this.sprite.dispose();
                }
                this.model = latestSnapshot.model;

                if (typeof this.templates[this.model] == 'undefined') {
                    console.error("cant load template for model "+this.model);
                    return;
                }
                this.sprite = this.templates[this.model].createInstance(this.id);
                this.sprite.isVisible = true;
            }

            if (this.scale != latestSnapshot.scale) {
                this.scale = latestSnapshot.scale;
                this.sprite.scaling = this.scale;
            }

            if (latestSnapshot.orientation) {

                var q = new BABYLON.Quaternion(latestSnapshot.orientation[1], latestSnapshot.orientation[2], latestSnapshot.orientation[3], latestSnapshot.orientation[0]);
                this.sprite.rotationQuaternion  = q
            }

            if (this.interpolationDelay <= 0) {
                this.sprite.position = latestSnapshot.position;
                return;
            }

            var fromSnapshot = this.getPreviousState(interpolationTime, latestSnapshot);
            if (fromSnapshot === false) {
                this.sprite.position = latestSnapshot.position;
                return;
            }

            coef = (interpolationTime - fromSnapshot.timestamp) / (latestSnapshot.timestamp - fromSnapshot.timestamp);
            if (coef < 0 || coef > 1) {
                this.sprite.position = latestSnapshot.position;
                return;
            }

            this.sprite.position = this.getInterpolated(fromSnapshot, latestSnapshot, coef);

            this.deleteOldSnapshots(fromSnapshot.timestamp);
        };

        /**
         * Delete all timestamps older than passed in timestamp
         *
         * @param timestamp
         */
        this.deleteOldSnapshots = function (timestamp) {
            if (typeof timestamp === 'undefined') {
                return;
            }
            // delete older than fromSnapshot
            for (var key in this.snapshots) {
                if (this.snapshots[key].timestamp < timestamp) {
                    this.snapshots.splice(key, 1);
                }
            }
        };

        /**
         *
         * @param timestamp
         * @returns {*}
         */
        this.getLatestState = function (timestamp) {
            var latestSnapshot;
            for (var key in this.snapshots) {
                if (timestamp >= this.snapshots[key].timestamp) {
                    latestSnapshot = this.snapshots[key];
                    break;
                }
            }
            if (typeof latestSnapshot === 'undefined') {
                return false;
            }
            return latestSnapshot;
        };

        /**
         *
         * @param timestamp
         * @param toState
         * @returns {*}
         */
        this.getPreviousState = function (timestamp, toState) {
            var fromSnapshot;
            if (typeof toState == 'undefined' || toState === false) {
                return false;
            }
            for (var key in this.snapshots) {
                if (timestamp > this.snapshots[key].timestamp) {
                    fromSnapshot = this.snapshots[key];
                    break;
                }
            }
            if (typeof fromSnapshot == 'undefined') {
                return false;
            }
            return fromSnapshot;
        };

        /**
         *
         * @param from
         * @param to
         * @param coef
         * @returns {{x: number, y: number}}
         */
        this.getInterpolated = function (from, to, coef) {
            var position = {x: 0, y: 0, z: 0},
                diffX,
                diffY,
                diffZ;

            diffX = to.position.x - from.position.x;
            if (Math.abs(diffX) < 0.1) {
                position.x = from.position.x;
            } else {
                position.x = from.position.x + coef * diffX;
            }
            diffY = to.position.y - from.position.y;
            if (Math.abs(diffY) < 0.1) {
                position.y = from.position.y;
            } else {
                position.y = from.position.y + coef * diffY;
            }

            diffZ = to.position.z - from.position.z;
            if (Math.abs(diffZ) < 0.1) {
                position.z = from.position.z;
            } else {
                position.z = from.position.z + coef * diffY;
            }
            return position;
        };


        /**
         *
         * @returns {PIXI.Sprite|*}
         */
        this.getSprite = function () {
            return this.sprite;
        };
    };

    var entity = {};

    entity.BUNNY = 2;

    entity.create = function (id, interpolationDelay, scene, templates) {
        var go;

        if (typeof interpolationDelay == 'undefined') {
            interpolationDelay = 0;
        }

        go = new GameObject(id, scene, templates);
        go.interpolationDelay = interpolationDelay;
        return go;
    };

    return entity;
});



