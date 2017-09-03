const assert = require('assert');
const publishers = require('../publishers')
const uuidv4 = require('uuid/v4');
const redis = require('redis');

describe('MockPublisher', () => {
    describe('#Publish', () => {
        it('should queue the published message', (done) => {
            var mockPub = publishers.FromURI(null, "mock://publish");
            mockPub.Publish("message").then(() => {
                assert.equal(mockPub.messages.length, 1);
                assert.equal(mockPub.messages[0], "message");
                done();
            });
        });
    });
    describe("Queue / Flush", () => {
        it('should queue the messages, then flush them', (done) => {
            var mockPub = publishers.FromURI(null, "mock://resumption");
            mockPub.QueueMessage("message 1").then(() => {
                return mockPub.QueueMessage("message 2")
            }).then(() => {
                return mockPub.QueueMessage("message 3")
            }).then(() => {
                assert.equal(mockPub.messages.length, 0);
                assert.equal(mockPub.queued.length, 3);
                return mockPub.FlushQueue();
            }).then(() => {
                assert.equal(mockPub.queued.length, 0);
                assert.equal(mockPub.messages.length, 3);
                assert.equal(mockPub.messages[0], "message 1")
                assert.equal(mockPub.messages[1], "message 2")
                assert.equal(mockPub.messages[2], "message 3")
                done();
            });
        });
    });
});
describe('RedisQueuePublisher', () => {
    if(!process.env.REDIS_URL) {
        console.log("No redis URL. Skipping redis tests");
        return;
    }
    var redisClient;
    var channelId;
    var queuePublisher;
    beforeEach(() => {
        [host, port] = process.env.REDIS_URL.split(":");
        redisClient = redis.createClient({host: host, port: port});
        channelId = uuidv4();
        queuePublisher = publishers.FromURI(redisClient, "queue://"+channelId);
    });
    afterEach(() => {
        redisClient.end(true);
    });
    describe('Publish', () => {
        it('should publish an item', (done) => {
            queuePublisher.Publish("message").then(() => {
                redisClient.llen(channelId, (err, data) => {
                    assert.equal(data, 1);
                    done();
                });
            });
        });
    });
    describe('Queue / Flush', () => {
        it('should queue then flush an item', (done) => {
            queuePublisher.QueueMessage("message").then(() => {
                return new Promise((resolve, reject) => {
                    redisClient.llen(channelId, (err, data) => {
                        assert.equal(data, 0);
                        resolve();
                    });
                });
            }).then(() => {
                return queuePublisher.FlushQueue();
            }).then(() => {
                redisClient.llen(channelId, (err, data) => {
                    assert.equal(data, 1);
                    done();
                });
            });
        });
    });
});
describe('RedisTopicPublisher', () => {
    if(!process.env.REDIS_URL) {
        console.log("No redis URL. Skipping redis tests");
        return;
    }
    var redisClient;
    var subClient;
    var channelId;
    var queuePublisher;
    var itemList;
    var subReady;
    beforeEach(() => {
        [host, port] = process.env.REDIS_URL.split(":");
        redisClient = redis.createClient({host: host, port: port});
        subClient = redis.createClient({host: host, port: port});
        channelId = uuidv4();
        topicPublisher = publishers.FromURI(redisClient, "topic://"+channelId);
        itemList = [];
        subClient.subscribe(channelId);
        subReady = new Promise((resolve, reject) => {
            subClient.on("subscribe", resolve);
        })
        subClient.on("message", (channel, message) => {
            itemList.push(message);
        });
    });
    afterEach(() => {
        redisClient.end(true);
        subClient.end(true);
    });
    describe('Publish', () => {
        it('should publish an item', (done) => {
            subReady.then(() => {
                return topicPublisher.Publish("message")
            }).then((err, data) => {
                setTimeout(() => {
                    assert.equal(itemList.length, 1);
                    done();
                }, 100);
            }).catch(console.log);
        });
    });
    describe('Queue / Flush', () => {
        it('should queue then flush an item', (done) => {
            a = 0;
            subReady.then(() => {
                return topicPublisher.QueueMessage("message");
            }).then(() => {
                assert.equal(itemList.length, 0);
            }).then(() => {
                return topicPublisher.FlushQueue();
            }).then(() => {
                setTimeout(() => {
                    assert.equal(itemList.length, 1);
                    done();
                }, 100);
            });
        });
    });
});
