define(["pixi"], function (pixi) {

    const ENTITY_WORLD = 1;
    const ENTITY_BUNNY = 2;

    var GameObject = function (texture) {

        this.texture = pixi.Texture.fromImage(texture);

        this.sprite = new pixi.Sprite(this.texture);

        this.sprite.anchor = {x: 0.5, y: 0.5};

        this.lastUpdate = null;

        this.server = new Array();

        /**
         * This updates on the rate that the server can push (50ms)
         *
         * @param message
         */
        this.serverUpdate = function (message) {
            this.server.unshift(message);
            while(this.server.length > 3) {
                this.server.pop();
            }
        };

        /**
         * this method is called approx every 50ms
         *
         * @param dTime - ms since last update
         */
        this.applyUpdates = function () {
//            if(this.lastUpdate === null) {
//                return;
//            }
//            this.server.push(this.lastUpdate);
//            this.lastUpdate = null;
//            // keep two server updates
//            while(this.server.length > 3) {
//                this.server.shift();
//            }
        }

        /**
         * this method is called approx every 16ms
         *
         * @param mSec
         */
        this.update = function(mSec) {
            if(this.server.length < 3) {
                return;
            }
//            this.sprite.position.x += (mSec) * this.velocity.x;
//            this.sprite.position.y += (mSec) * this.velocity.y;
              //this.sprite.position = this.server[this.server.length-1].position;
        }

        /**
         * this method is called every approx 16ms
         *
         * @param mSec
         */
        this.interpolate = function (range) {
            if(this.server.length < 3) {
                return;
            }
            if(range <= 0 || range >= 1) {
                return;
            }

            var diffX =  this.server[0].position.x - this.sprite.position.x;
            if(Math.abs(diffX) < 2) {
                this.sprite.position.x = this.server[0].position.x;
            } else {
                this.sprite.position.x += diffX * range;
            }

            var diffY =  this.server[0].position.y - this.sprite.position.y;
            if(Math.abs(diffY) < 2) {
                this.sprite.position.y = this.server[0].position.y;
            } else {
                this.sprite.position.y += diffY * range;
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



