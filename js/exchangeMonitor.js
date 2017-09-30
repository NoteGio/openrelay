const redis = require("redis");
const ZeroEx = require("0x.js")
const monitor = require("./monitor");
const Exchange = artifacts.require("Exchange");

module.exports = function(done){
    var redisClient = redis.createClient(process.argv[4]);
    var notifyURL = process.argv[5];
    // TODO: Once we have affiliates, this will need to come from redis
    var feeRecipients = [process.argv[6]];
    zeroEx = new ZeroEx.ZeroEx(web3.currentProvider);
    zeroEx.exchange.getContractAddressAsync().then((contractAddress) => {
        return Exchange.at(contractAddress);
    }).then((exchangeContract) => {
        for(var feeRecipient of feeRecipients) {
            var fillFilterCreator = (options) => {
                return exchangeContract.LogFill(options, {feeRecipient: feeRecipient})
            }
            var fillTransform = (data) => {
                return {
                    orderHash: data.args.orderHash,
                    filledTakerTokenAmount: data.args.filledTakerTokenAmount,
                }
            }
            var cancelFilterCreator = (options) => {
                return exchangeContract.LogCancel(options, {feeRecipient: feeRecipient})
            }
            var cancelTransform = (data) => {
                return {
                    orderHash: data.args.orderHash,
                    cancelledTakerTokenAmount: data.args.cancelledTakerTokenAmount,
                }
            }
            monitor(redisClient, notifyURL, fillFilterCreator, web3, fillTransform);
            monitor(redisClient, notifyURL, cancelFilterCreator, web3, cancelTransform);
        }
    });
};
