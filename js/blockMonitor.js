const redis = require("redis");
const monitor = require("./monitor");

module.exports = function(done){
    var redisClient = redis.createClient(process.argv[4]);
    var notifyURL = process.argv[5];
    var filterCreator = function() { return web3.eth.filter("latest"); }
    monitor(redisClient, notifyURL, filterCreator, web3, (data) => {
        return new Promise((resolve, reject) => {
            web3.eth.getBlock(data, function(error, result) {
                if(error){
                    reject(error);
                } else {
                    if (!result) {
                        reject("Result is empty");
                        return;
                    }
                    var blockObj = { "hash": result['hash'], "number": result['number'], "bloom": result['logsBloom'] };
                    resolve(blockObj);
                }
            })
        });
    });
};
