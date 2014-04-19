define(["pixi"], function (pixi){

    var my = {};

    my.create = function(x,y) {
        // create a texture from an image path
        var texture = pixi.Texture.fromImage("sprites/bunny.png");

        // create a new Sprite using the texture
        var bunny = new pixi.Sprite(texture);

        // center the sprites anchor point
        bunny.anchor = {x: 0.5, y: 0.5}

        // move the sprite t the center of the screen
        bunny.position = { x: x, y: y};
        return bunny;
    }

    return my;

});



