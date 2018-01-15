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
            if(err) {
                console.log(err);
                var stack = new Error().stack
                console.log( stack )
                // If we get an error, exit. Docker should restart the process,
                // and it can go through the resumption process.
                //
                // We may need to recover more gracefully, but we'll try this
                // for now.
                process.exit(1);
            }
            if(resuming){
                clearTimeout(resumptionTimeout);
                resumptionTimeout = setTimeout(function(){
                    web3.eth.getBlockNumber((err, currentBlock) => {
                        resuming = false;
                        channel.FlushQueue();
                        redisClient.set(blockKey, currentBlock);
                        lastBlockNumber = currentBlock;
                    })
                }, 5000);
                channel.QueueMessage(JSON.stringify(transform(data)));
            } else {
                web3.eth.getBlockNumber((err, currentBlock) => {
                    var payload = JSON.stringify(transform(data));
                    console.log(`Block ${currentBlock} - published '${payload}'`);
                    channel.Publish(payload);
                    if(lastBlockNumber != currentBlock && currentBlock != null) {
                        redisClient.set(blockKey, currentBlock);
                        lastBlockNumber = currentBlock;
                    }
                });
            }

        });
    }).catch((err) => {
        console.log("Failed to get initial block number: " + err);
    })
}
