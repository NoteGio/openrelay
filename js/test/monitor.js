const monitor = require("../monitor");
const assert = require('assert');
const publishers = require('../publishers')

function MockRedisClient() {
    var data = {};
    this.get = (key, callback) => {
        callback(null, data[key] || null);
    }
    this.set = (key, value) => {
        data[key] = value;
    }
}

function MockWeb3(blockNumber) {
    this.eth = {blockNumber: blockNumber || 0}
    this.addBlock = () => {
        this.eth.blockNumber++;
    }
}

function MockFilter() {
    watchers = [];
    this.watch = (cb) => {
        watchers.push(cb)
    }
    this.trigger = (err, data) => {
        for(var i = 0; i < watchers.length; i++) {
            watchers[i](err, data);
        }
    }
}

function tick() {
    return new Promise((resolve, reject) => {
        setTimeout(resolve);
    })
}

function passthrough(data) {
    return data;
}

describe('MockMonitor', () => {
    describe('Test message filter', () => {
        var redisClient;
        var web3;
        beforeEach(() => {
            redisClient = new MockRedisClient();
            web3 = new MockWeb3();
        });
        it("should start from a prior block and flush after a timeout", (done) => {
            redisClient.set("mock://testresumption::blocknumber", 1);
            web3.eth.blockNumber = 1;
            var filter = new MockFilter();
            var channel = publishers.FromURI(null, "mock://testresumption");
            monitor(redisClient, "mock://testresumption", () => {return filter}, web3, passthrough).then(() => {
                filter.trigger(null, "message");
                return tick();
            }).then(() => {
                assert.equal(channel.messages.length, 0);
                assert.equal(channel.queued.length, 1);
            }).then(() => {
                return new Promise((resolve, reject) => {
                    setTimeout(() => {
                        assert.equal(channel.messages.length, 1);
                        assert.equal(channel.queued.length, 0);
                        resolve();
                    }, 5001)
                });
            }).then(() => {
                filter.trigger(null, "message2");
                return tick();
            }).then(() => {
                assert.equal(channel.messages.length, 2);
                assert.equal(channel.queued.length, 0);
                done();
            });
        }).timeout(6000);
        it("should start monitoring straight away, with no resumption period", (done)=> {
            web3.eth.blockNumber = 1;
            var filter = new MockFilter();
            var channel = publishers.FromURI(null, "mock://testinitialize");
            monitor(redisClient, "mock://testinitialize", () => {return filter}, web3, passthrough).then(() => {
                filter.trigger(null, "message");
                return tick();
            }).then(() => {
                assert.equal(channel.messages.length, 1);
                assert.equal(channel.queued.length, 0);
                done();
            })
        });
    });
});
