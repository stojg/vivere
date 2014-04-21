require(["screen", "websocket", 'pixi', 'entity', "ui"], function(screen, websocket, pixi, entity, ui) {

    var commands = new Array();

    // create an new instance of a pixi stage
    var stage = new pixi.Stage(0x666666);

    ui.setMessage("Loading...", stage)

    var entities = new Array();

    var conn = websocket.connect(commands);

    websocket.afterConnect(conn, function(){
        ui.clearMessage(stage);

        conn.send(JSON.stringify({
            Event: "World",
            Message: "getState"
        }));

        var renderer = false;

        requestAnimFrame(animate);

        function animate() {
            while (commands.length > 0) {
                var command = commands.pop();
                var message = JSON.parse(command);
                if(message.Event === "Entity") {
                    updateEntity(message.Message);
                } else if(message.Event === "World") {
                    // create a renderer instance.
                    renderer = pixi.autoDetectRenderer(message.Message.Width,  message.Message.Height);
                    // add the renderer view element to the DOM
                    document.body.appendChild(renderer.view);
                } else {
                    console.log(message);
                }
            }
            requestAnimFrame(animate);
            if(renderer != false) {
                renderer.render(stage);
            }

        }
    })

    var updateEntity = function(entityMsg) {
        if(!entities[entityMsg.Id]) {
            entities[entityMsg.Id] = entity.create();
            stage.addChild(entities[entityMsg.Id]);
        }
        entities[entityMsg.Id].rotation = entityMsg.Rotation
        entities[entityMsg.Id].position.x = entityMsg.Position.X
        entities[entityMsg.Id].position.y = entityMsg.Position.Y
    }
});