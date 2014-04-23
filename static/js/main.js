require(["screen", "websocket", 'pixi', 'entity', "ui", "commands"], function(screen, websocket, pixi, entity, ui, commands) {

    var framesPerSecond = 1000 / 60;
    // create an new instance of a pixi stage
    var stage = new pixi.Stage(0x666666);

    var entities = new Array();

    var connected = false;

    var gameTick = 0;
    var renderLoopInterval = null;

    var renderer = pixi.autoDetectRenderer(1000, 600);
    document.body.appendChild(renderer.view);

    websocket.connect(function() {
        connected = true;
        frame();
    }, getState);

    function sendCmd() {
        if(commands.get() == 0) {
            return false;
        }
        var cmd = new DataStream();
        cmd.writeUint32(commands.get());
        connected = websocket.send(cmd.buffer);
    }

    var current = Date.now();

    function frame() {
        renderLoopInterval = setTimeout(function() {

            var now = Date.now();
            var elapsed = now - current;
            current = now;

            if(!connected) {
                console.log('Server died, please reload page!');
                clearInterval(renderLoopInterval);
                return;
            }

            window.requestAnimationFrame(frame);

            for (key in entities) {
                // @todo some clever lerp:ing
            }

            renderer.render(stage);
            sendCmd();
        }, framesPerSecond);
    }

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
                entities[id] = entity.create(modelId);
                stage.addChild(entities[id]);
            }

            // rotation
            if ((bitMask & (1<<1))>0) {
                entities[id].rotation = buf.readFloat32();
            }

            // angular velocity
            if ((bitMask & (1<<2))>0) {
                entities[id].angularVel = buf.readFloat32();
            }

            // pos
            if ((bitMask & (1<<3))>0) {
                var pos = buf.readFloat64Array(2);
                entities[id].position.x = pos[0];
                entities[id].position.y = pos[1];
            }

            // vel
            if ((bitMask & (1<<4))>0) {
                var vel = buf.readFloat64Array(2);
                entities[id].velocity = {x: vel[0], y: vel[1]};
            }

            // size
            if ((bitMask & (1<<5))>0) {
                var size = buf.readFloat64Array(2);
                entities[id].size = {x: size[0], y: size[1]};
            }
        }
    }
});