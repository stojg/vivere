require(["screen", "websocket", 'pixi', 'entity', "ui"], function(screen, websocket, pixi, entity, ui) {

    var commands = new Array();

    // create an new instance of a pixi stage
    var stage = new pixi.Stage(0x666666);

    // create a renderer instance.
    var renderer = pixi.autoDetectRenderer(screen.size.width,  screen.size.height);

    ui.setMessage("Loading...", stage);

    // add the renderer view element to the DOM
    document.body.appendChild(renderer.view);

    var entities = new Array();

    var conn = websocket.connect(commands);

    websocket.afterConnect(conn, function(){
        ui.clearMessage(stage);

        requestAnimFrame(animate);
        function animate() {
            while (commands.length > 0) {
                var command = commands.pop();
                var message = JSON.parse(command);
                if(!entities[message.Id]) {
                    console.log(message);
                    entities[message.Id] = entity.create();
                    stage.addChild(entities[message.Id]);
                }
                entities[message.Id].rotation = message.Rotation
                entities[message.Id].position.x = message.Position.X
                entities[message.Id].position.y = message.Position.Y
                //console.log(JSON.parse(command));
            }
            requestAnimFrame(animate);
            renderer.render(stage);
        }
    })
});