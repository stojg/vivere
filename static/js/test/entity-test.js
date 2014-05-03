define(['src/entity', 'lib/chai'], function (entity, chai) {

    var assert = chai.assert;

    describe('For an entity', function () {

        var latestState = {
            "id": 1,
            "timestamp": 10,
            "tick": 126,
            "rotation": 0,
            "position": { "x": 2, "y": 4 },
            "velocity": { "x": 0, "y": 0 },
            "size": { "x": 10, "y": 20 },
            "state": 0
        };

        var previousState = {
            "id": 2,
            "timestamp": 20,
            "tick": 127,
            "rotation": 0,
            "position": { "x": 4, "y": 8 },
            "velocity": { "x": 0, "y": 0 },
            "size": { "x": 10, "y": 20 },
            "state": 0
        };

//        it('get interpolated', function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            var from = { position: {x: 0, y: 0 } };
//            var to = { position: {x: 10, y: 10 } };
//            var result = bunny.getInterpolated(from, to, 0.5);
//            assert.equal(result.x, 5, 'position.x');
//            assert.equal(result.y, 5, 'position.y');
//        });
//
//        it('should be repositioned after a server update', function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            bunny.serverUpdate(latestState);
//            bunny.update(10);
//            assert.equal(bunny.getSprite().position.x, 2);
//            assert.equal(bunny.getSprite().position.y, 4);
//        });
//
//        it("there should not be a valid to_snapshot when there is no server updates", function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            //bunny.applyServerUpdates();
//            var to = bunny.getLatestState(1);
//            assert.equal(to, false);
//        });
//
//        it("there should be no to_snapshot with a timestamp ", function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            bunny.serverUpdate(latestState);
//            bunny.applyServerUpdates();
//            var to = bunny.getLatestState(0);
//            assert.equal(to, false);
//        });
//
//        it("there should be a to_snapshot", function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            bunny.serverUpdate(latestState);
//            bunny.applyServerUpdates();
//            var to = bunny.getLatestState(10);
//            assert.equal(to, latestState);
//        });
//
//        it("there should not be a previous state after one update", function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            bunny.serverUpdate(latestState);
//            bunny.applyServerUpdates();
//            var to = bunny.getLatestState(5);
//            var from = bunny.getPreviousState(5, to);
//            assert.equal(from, false);
//        });
//
//        it("there should not be a previous state because the delay is to small", function () {
//            var bunny = entity.create(entity.BUNNY, 0);
//            bunny.serverUpdate(latestState);
//            bunny.serverUpdate(previousState);
//            bunny.applyServerUpdates();
//            var to = bunny.getLatestState(5);
//            var from = bunny.getPreviousState(6, to);
//            assert.equal(from, false);
//        });

        it("there should be a previous state", function () {
            var bunny = entity.create(entity.BUNNY, 0);
            bunny.serverUpdate(latestState);
            bunny.serverUpdate(previousState);
            bunny.applyServerUpdates();
            var to = bunny.getLatestState(10);
            assert.equal(to, latestState, 'to state should exists');
            var from = bunny.getPreviousState(21, to);
            assert.equal(from, previousState);
        });

    });
});