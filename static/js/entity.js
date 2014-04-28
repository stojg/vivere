define(["pixi"], function (pixi) {

    const ENTITY_WORLD = 1;
    const ENTITY_BUNNY = 2;

    var GameObject = function (texture) {

        this.texture = pixi.Texture.fromImage(texture);

        this.sprite = new pixi.Sprite(this.texture);

        this.sprite.anchor = {x: 0.5, y: 0.5};

        this.updates = new Array();

        /**
         *
         * @param message
         */
        this.serverUpdate = function(message) {
            this.updates.push(message)
        };

        /**
         *
         * @param dTime - ms since last update
         */
        this.update = function(dTime) {
            var mess = this.updates.pop();
            if(!mess) {
                return;
            }
            this.sprite.position = mess.position;
            this.sprite.rotation = mess.rotation;
            //console.log('update');
            //console.log(dTime);
        }

        this.getSprite = function() {
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



