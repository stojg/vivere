define(["pixi", "gamestate"], function (pixi, gamestate) {

    const ENTITY_WORLD = 1;
    const ENTITY_BUNNY = 2;

    var GameObject = function (texture) {

        this.texture = pixi.Texture.fromImage(texture);

        this.sprite = new pixi.Sprite(this.texture);

        this.sprite.anchor = {x: 0.5, y: 0.5};

        this.interpolationDelay = 100;

        /**
         * Contains a list of queued updates from the server
         *
         * Older update -> last update
         *
         * @type {Array}
         */
        this.server = new Array();

        /**
         * Contains a list of updates in the past, used for interpolation
         * Last snapshot -> older
         *
         * @type {Array}
         */
        this.snapshots = new Array();

        /**
         * serverUpdate is called when the server sends an update command
         * to this entity
         *
         * @param message
         */
        this.serverUpdate = function (message) {
            this.server.push(message);
        };

        /**
         *
         * @param mSec
         */
        this.update = function(tFrame) {
            // Move queued server updates to the snapshot array
            var msg = this.server.pop();
            while(typeof msg !== 'undefined') {
                this.snapshots.unshift(msg);
                msg = this.server.pop();
            }

            if(this.interpolationDelay < 1) {
                this.sprite.position.x = this.snapshots[0].position.x;
                this.sprite.position.y = this.snapshots[0].position.y;
                this.deleteOldSnapshots(this.snapshots[0].timestamp);
            } else {
                this.interpolate(tFrame);
            }
        }

        /**
         *
         * @param tFrame
         */
        this.interpolate = function (tFrame) {
            var interpolationTime = tFrame - this.interpolationDelay,
                latestSnapshot,
                fromSnapshot,
                coef,
                diffX,
                diffY;

            for(key in this.snapshots) {
                if(tFrame > this.snapshots[key].timestamp) {
                    latestSnapshot = this.snapshots[key];
                    break;
                }
            }

            if(typeof latestSnapshot === 'undefined') {
                console.info("There is no snapshot older then tFrame");
                return;
            }

            for(key in this.snapshots) {
                if(interpolationTime > this.snapshots[key].timestamp) {
                    fromSnapshot = this.snapshots[key];
                    break;
                }
            }

            if(typeof fromSnapshot === 'undefined') {
                console.info("Not enough snapshots to interpolate between");
                this.sprite.position.x = latestSnapshot.position.x;
                this.sprite.position.y = latestSnapshot.position.y;
                return;
            }

            this.deleteOldSnapshots(fromSnapshot.timestamp);

            if(latestSnapshot.timestamp == fromSnapshot.timestamp) {
                console.info("Latest and from snapshots are the same, interpolationDelay too low.");
                this.sprite.position.x = latestSnapshot.position.x;
                this.sprite.position.y = latestSnapshot.position.y;
                return;
            }

            if(latestSnapshot.timestamp < fromSnapshot.timestamp) {
               console.error("From snapshot have a newer timestamp than latest snapshot");
               return;
            }

            coef = (interpolationTime  - fromSnapshot.timestamp) / (latestSnapshot.timestamp - fromSnapshot.timestamp);

            if(coef < 0 || coef > 1) {
                console.error("Interpolation coef is out of bounds: " + coef);
                this.sprite.position.x = latestSnapshot.position.x;
                this.sprite.position.y = latestSnapshot.position.y;
                return;
            }

            diffX = latestSnapshot.position.x - fromSnapshot.position.x;
            if(Math.abs(diffX) < 0.1 ) {
                this.sprite.position.x = latestSnapshot.position.x;
            } else {
                this.sprite.position.x = fromSnapshot.position.x + coef * diffX;
            }

            diffY = latestSnapshot.position.y - fromSnapshot.position.y;
            if(Math.abs(diffY) < 0.1 ) {
                this.sprite.position.y = latestSnapshot.position.y;
            } else {
                this.sprite.position.y = fromSnapshot.position.y + coef * diffY;
            }
        }

        /**
         * Delete all timestamps older than passed in timestamp
         *
         * @param timestamp
         */
        this.deleteOldSnapshots = function(timestamp) {
            if(typeof timestamp === 'undefined') {
                console.error("timestamp undefined");
                return;
            }
            // delete older than fromSnapshot
            for(key in this.snapshots) {
                if(this.snapshots[key].timestamp < timestamp) {
                    delete(this.snapshots[key]);
                }
            }
        }

        /**
         *
         * @returns {PIXI.Sprite|*}
         */
        this.getSprite = function () {
            return this.sprite;
        }
    }

    var entity = {};

    entity.create = function (type) {

        if (type === ENTITY_BUNNY) {
            return new GameObject("sprites/bunny.png");
        }
        if (type === ENTITY_WORLD) {
            var entity = new pixi.Stage();
            entity.anchor = {x: 0.5, y: 0.5};
            return entity;
        }

        throw new Error("Tried to create a model without an exiting type '" + type + "'");
    }
    return entity;
});



