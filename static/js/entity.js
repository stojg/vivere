define(["pixi"], function (pixi){

    var my = {};

    my.create = function() {
        // create a texture from an image path
        var texture = pixi.Texture.fromImage("sprites/bunny.png");
        // create a new Sprite using the texture
        var bunny = new pixi.Sprite(texture);
        // center the sprites anchor point
        bunny.anchor = {x: 0.5, y: 0.5}
        return bunny;
    }
    return my;
});



