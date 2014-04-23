define(["pixi"], function (pixi){

    var my = {};

    const ENTITY_WORLD = 1;
    const ENTITY_BUNNY = 2;

    my.create = function(type) {

        if(type === ENTITY_BUNNY ) {
            // create a texture from an image path
            var texture = pixi.Texture.fromImage("sprites/bunny.png");
            // create a new Sprite using the texture
            var bunny = new pixi.Sprite(texture);
            // center the sprites anchor point
            bunny.anchor = {x: 0.5, y: 0.5}
            return bunny;
        }
        if(type === ENTITY_WORLD ) {
            var entity = new pixi.Stage();
            entity.anchor = {x: 0.5, y: 0.5};
            return entity;
        }

        throw new Error("Tried to create a model without an exiting type '" +type+"'");
    }
    return my;
});



