const IssueTokenFactory = artifacts.require("IssueTokenFactory");
const IssueToken = artifacts.require("IssueToken");
const redis = require("redis");

function issueToken(redisName, tokenFactory, redisClient) {
    var tokenAddress;
    tokenFactory.newToken.call().then((address) => {
        tokenAddress = address;
        redisClient.set(redisName + "::address", address.substr(2))
    })
    return tokenFactory.newToken().then(() => {
        return IssueToken.at(tokenAddress).issue(10**27);
    }).then(() => {
        return IssueToken.at(tokenAddress).transfer(web3.eth.accounts[1], 10**25);
    });
}

module.exports = function(done){
    var redisClient = redis.createClient(process.argv[4]);
    IssueTokenFactory.deployed().then((tokenFactory) => {
        return issueToken("feeToken", tokenFactory, redisClient).then(() => {
            return issueToken("tokenX", tokenFactory, redisClient);
        }).then(() => {
            return issueToken("tokenY", tokenFactory, redisClient);
        }).then(() => {
            return issueToken("tokenZ", tokenFactory, redisClient);
        }).then(done)
    })
}
