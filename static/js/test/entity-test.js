define(['src/entity', 'lib/chai'], function (entity, chai) {

    describe('Entity', function () {

        it('sholdu be able to recieve a an update from the server', function () {
            var bunny = entity.create(entity.BUNNY);
            bunny.serverUpdate({'test': 'ss'});
        });

    });
});