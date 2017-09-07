const redis = require("redis");
const monitor = require("./monitor");

module.exports = function(done){
    var redisClient = redis.createClient(process.argv[4]);
    var notifyURL = process.argv[5];
    var filterCreator = function() { return web3.eth.filter("latest"); }
    monitor(redisClient, notifyURL, filterCreator, web3, (data) => {return data;});
};
