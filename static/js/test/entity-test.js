define(['src/entity', 'lib/chai'], function (entity, chai) {

    var assert = chai.assert

    describe('Entity', function () {

        var message = {
            "id": 1,
            "timestamp": 220,
            "tick": 126,
            "rotation": 0,
            "position": { "x": 2, "y": 4 },
            "velocity": { "x": 0, "y": 0 },
            "size": { "x": 10, "y": 20 },
            "state": 0
        };

        it('should be repositioned after a server update', function () {
            var bunny = entity.create(entity.BUNNY, 0);
            bunny.serverUpdate(message);
            bunny.update(1);
            assert.equal(bunny.getSprite().position.x, 2);
            assert.equal(bunny.getSprite().position.y, 4);
        });

    });
});