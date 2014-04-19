require(["screen", "websocket", 'pixi', 'entity', "ui"], function(screen, websocket, pixi, entity, ui) {

    var commands = new Array();

    // create an new instance of a pixi stage
    var stage = new pixi.Stage(0x666666);

    // create a renderer instance.
    var renderer = pixi.autoDetectRenderer(screen.size.width,  screen.size.height);

    ui.setMessage("Loading...", stage);

    // add the renderer view element to the DOM
    document.body.appendChild(renderer.view);

    var bunny = entity.create(screen.size.width/2,  screen.size.height/2);

    var conn = websocket.connect(commands);

    websocket.afterConnect(conn, function(){
        ui.clearMessage(stage);
        stage.addChild(bunny);
        requestAnimFrame(animate);
        function animate() {
            while (commands.length > 0) {
                var command = commands.pop();
                bunny.rotation = JSON.parse(command).Rotation
            }
            requestAnimFrame(animate);
            renderer.render(stage);
        }
    })
});