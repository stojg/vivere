define(["pixi"], function (pixi){

    var my = {};

    my.message;

    my.setMessage = function(message, stage) {
        this.message = new pixi.Text(message);
        this.message.setStyle({
            fill: "#fff"
        })
        stage.addChild(this.message);
    }

    my.clearMessage = function(stage) {
       stage.removeChild(this.message);
    }
    return my;
});
