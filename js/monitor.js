var redis = require("redis");

/*
This function needs to take a filter. When events come in on the filter,
publish them to the queue via the Redis URL. Every time we publish an event,
if the block number has changed since the previous block we record the block
number in redis.

When we first start up, if the process we need to start in "resumption mode".
In resumption mode, our filter should go from the block stored in redis to
"latest". Resumption mode must queue up all of the messages it processes until
resumption mode finishes. When resumption mode starts, a timeout should be set
while messages are being processed. Each time a message is processed, it should
reset the timeout. When the timeout executes we can assume that resumption
has finished, so the queue should be flushed, then the block number should be
updated.
*/

var publishers = require("./publishers");

// TODO: Make sure filter will be recreated if the backend RPC server changes

module.exports = function(redisClient, notificationChannel, filterCreator, web3, transform){
    var blockKey = notificationChannel + "::blocknumber";
    channel = publishers.FromURI(redisClient, notificationChannel);
    return new Promise((resolve, reject) => {
        redisClient.get(blockKey, function(err, data) {
            if(err){
                reject(err);
            } else {
                resolve(data);
            }
        })
    }).then((blockNumber) => {
        var resuming;
        var resumptionTimeout;
        if(blockNumber === null) {
            blockNumber = "latest";
            resuming = false;
        } else {
            // blockNumber from Redis is the last completed block, so we
            // want to resume from the next block.
            blockNumber = parseInt(blockNumber) + 1;
            resuming = true;
        }
        var watcher = filterCreator({fromBlock: blockNumber, toBlock: "latest"});
        lastBlockNumber = blockNumber;
        watcher.watch((err, data) => {
            currentBlock = web3.eth.blockNumber;
            if(resuming){
                clearTimeout(resumptionTimeout);
                resumptionTimeout = setTimeout(function(){
                    resuming = false;
                    channel.FlushQueue();
                    redisClient.set(blockKey, web3.eth.blockNumber);
                    lastBlockNumber = web3.eth.blockNumber;
                }, 5000);
                channel.QueueMessage(JSON.stringify(transform(data)));
            } else {
                channel.Publish(JSON.stringify(transform(data)));
                if(lastBlockNumber != currentBlock) {
                    redisClient.set(blockKey, currentBlock);
                    lastBlockNumber = currentBlock;
                }
            }

        });
    }).catch((err) => {
        console.log("Failed to get initial block number: " + err);
    })
}
