require(["screen", "websocket", 'pixi', 'entity', "ui", "commands"], function(screen, websocket, pixi, entity, ui, commands) {

    var framesPerSecond = 1000 / 60;
    // create an new instance of a pixi stage
    var stage = new pixi.Stage(0x666666);

    var entities = new Array();

    var connected = false;

    var renderLoopInterval = null;

    var renderer = pixi.autoDetectRenderer(1000, 600);
    document.body.appendChild(renderer.view);

    websocket.connect(function() {
        connected = true;
        frame();
    }, getState);

    function sendCmd() {
        var cmd = new DataStream();
        cmd.writeUint32(commands.get());
        connected = websocket.send(cmd.buffer);
    }

    function frame() {
        renderLoopInterval = setTimeout(function() {
            if(!connected) {
                console.log('Server died, please reload page!');
                clearInterval(renderLoopInterval);
                return;
            }
            window.requestAnimationFrame(frame);
            renderer.render(stage);
            sendCmd();
        }, framesPerSecond);
    }

    var ents = [];

    function getState(evt) {
        var buf = new DataStream(evt.data)

        // Number of entities
        var nEnts =  buf.readUint16();

        for(i = 0; i < nEnts; i++) {

            // get the bitmask
            var bitMask = buf.readUint8();

            // id
            var id = buf.readUint16();

            // model
            if ((bitMask & (1<<0))>0) {
                var modelId = buf.readUint16();
                ents[id] = entity.create(modelId);
                stage.addChild(ents[id]);
            }

            // rotation
            if ((bitMask & (1<<1))>0) {
                ents[id].rotation = buf.readFloat32();
            }

            // angular velocity
            if ((bitMask & (1<<2))>0) {
                ents[id].angularVel = buf.readFloat32();
            }

            // pos
            if ((bitMask & (1<<3))>0) {
                var pos = buf.readFloat64Array(2);
                ents[id].position.x = pos[0];
                ents[id].position.y = pos[1];
            }

            // vel
            if ((bitMask & (1<<4))>0) {
                var vel = buf.readFloat64Array(2);
                ents[id].velocity = {x: vel[0], y: vel[1]};
            }

            // size
            if ((bitMask & (1<<5))>0) {
                var size = buf.readFloat64Array(2);
                ents[id].size = {x: size[0], y: size[1]};
            }
        }
    }
});