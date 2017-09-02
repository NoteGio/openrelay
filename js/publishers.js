function QueuePublisher(redisClient, queueName) {
    var resumptionQueue = queueName + "rqueue";
    this.Publish = function(message) {
        return new Promise((resolve, reject) => {
            redisClient.lpush(queueName, message, resolve);
        });
    }
    this.QueueMessage = function(message) {
        return new Promise((resolve, reject) => {
            redisClient.lpush(resumptionQueue, message, resolve);
        });
    }
    this.FlushQueue = function() {
        return new Promise((resolve, reject) => {
            redisClient.llen(resumptionQueue, (err, data) => {
                if(err) {
                    reject(err);
                } else {
                    resolve(parseInt(data));
                }
            });
        }).then((count) => {
            batch = redisClient.batch();
            for(var i = 0; i < count; i++) {
                batch = batch.rpoplpush(resumptionQueue, queuName);
            }
            batch.exec();
        });
    }
}

function TopicPublisher(redisClient, topicName) {
    var resumptionQueue = topicName + "rqueue";
    this.Publish = function(message) {
        return new Promise((resolve, reject) => {
            redisClient.publish(topicName, message, resolve);
        })
    }
    this.QueueMessage = function(message) {
        return new Promise((resolve, reject) => {
            redisClient.lpush(resumptionQueue, message, resolve);
        })
    }
    this.FlushQueue = function() {
        return new Promise((resolve, reject) => {
            redisClient.llen(resumptionQueue, (err, data) => {
                if(err) {
                    reject(err);
                } else {
                    resolve(parseInt(data));
                }
            });
        }).then((count) => {
            batch = redisClient.batch();
            for(var i = 0; i < count; i++) {
                batch = batch.eval("return redis.call('PUBLISH', "+ topicName +", redis.call('RPOP', "+ resumptionQueue +"))");
            }
            batch.exec();
        });
    }
}

function MockPublisher(redisClient, name) {
    this.messages = [];
    this.queued = [];
    this.Publish = function(message) {
        return new Promise((resolve, reject) => {
            this.messages.push(message);
            resolve();
        });
    }
    this.QueueMessage = function(message) {
        return new Promise((resolve, reject) => {
            queued.push(message);
            resolve();
        });
    }
    this.FlushQueue = function() {
        return new Promise((resolve, reject) => {
            for(var i = 0; i < queued.length; i++) {
                messages.push(queued.pop());
            }
            resolve();
        });
    }
}

function FromURI(redisClient, channelUri) {
    if(channelUri.startsWith("topic://")) {
        return new TopicPublisher(redisClient, channelUri.substr("topic://".length));
    }
    if(channelUri.startsWith("queue://")) {
        return new QueuePublisher(redisClient, channelUri.substr("queue://".length));
    }
    if(channelUri.startsWith("mock://")) {
        return new MockPublisher();
    }
}

module.exports = {
    QueuePublisher: QueuePublisher,
    TopicPublisher: TopicPublisher,
    MockPublisher: MockPublisher,
    FromURI: FromURI
}
