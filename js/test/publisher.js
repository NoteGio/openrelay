var assert = require('assert');
var publishers = require('../publishers')

describe('MockPublisher', () => {
    describe('#Publish', () => {
        it('should queue the published message', (done) => {
            var mockPub = publishers.FromURI(null, "mock://");
            mockPub.Publish("message").then(() => {
                assert.equal(mockPub.messages.length, 1);
                assert.equal(mockPub.messages[0], "message");
                done();
            });
    });
    });
});
