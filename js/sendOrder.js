const Token = artifacts.require("Token");
const redis = require("redis");
const ZeroEx = require("0x.js")
const http = require("http");

function getAsync(redisClient, key) {
    return new Promise((resolve, reject) => {
        redisClient.get(key, (err, value) => {
            if(err){
                reject(err);
                return
            }
            resolve(value);
        });
    });
}

var requiredFee = new web3.BigNumber("500000000000000000");

module.exports = function(done){
    web3.eth.getAccounts((err, accounts) => {
        var defaultAccount = accounts[0];
        var redisClient = redis.createClient(process.argv[4]);
        var zeroEx = new ZeroEx.ZeroEx(web3.currentProvider);
        var order = {
            expirationUnixTimestampSec: new web3.BigNumber(Math.floor(Date.now()/1000) + (60*60*24)),
            feeRecipient: "0xc22d5b2951db72b44cfb8089bb8cd374a3c354ea",
            maker: defaultAccount,
            makerFee: requiredFee.div(2),
            makerTokenAmount: new web3.BigNumber(10000),
            salt: ZeroEx.ZeroEx.generatePseudoRandomSalt(),
            taker: "0x0000000000000000000000000000000000000000",
            takerFee: requiredFee.minus(requiredFee.div(2)), // account for rounding
            takerTokenAmount: new web3.BigNumber(10000),
        };
        zeroEx.proxy.getContractAddressAsync().then((proxyAddress) => {
            console.log(proxyAddress);
            return Promise.all([
                getAsync(redisClient, "feeToken::address").then((feeAddress) => {
                    console.log(feeAddress);
                    Token.at("0x" + feeAddress).approve(proxyAddress, requiredFee.div(2));
                }),
                getAsync(redisClient, "tokenX::address").then((address) => {
                    console.log(address);
                    order["makerTokenAddress"] = "0x" + address;
                    Token.at("0x" + address).approve(proxyAddress, new web3.BigNumber(10000));
                }),
                getAsync(redisClient, "tokenY::address").then((address) => {
                    console.log(address);
                    order["takerTokenAddress"] = "0x" + address;
                }),
                zeroEx.exchange.getContractAddressAsync().then((address) => {
                    console.log(address);
                    order["exchangeContractAddress"] = address;
                })
            ])
        }).then(() => {
            var orderHash = ZeroEx.ZeroEx.getOrderHashHex(order);
            return zeroEx.signOrderHashAsync(orderHash, defaultAccount);
        }).then((signature) => {
            order["ecSignature"] = signature;
            var start = Date.now();
            var req = http.request({
                host: "ingest",
                port: "8080",
                path: "/v0.0/order",
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                }
            }, (res) => {
                console.log(res.statusCode);
                res.on('data', console.log);
                res.on('end', () => {
                    console.log(Date.now() - start);
                    redisClient.quit();
                    done();
                });
            });
            req.write(JSON.stringify(order));
            req.end();
        }).catch(console.log);
    });
}
