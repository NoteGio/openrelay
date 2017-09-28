const ZeroEx = require("0x.js");
const fs = require("fs");
const util = require("ethereumjs-util");
const bignumber = require("bignumber.js");
const Exchange = artifacts.require("Exchange");

module.exports = function(done){
    order = fixBn(JSON.parse(fs.readFileSync(process.argv[4])));
    zeroEx = new ZeroEx.ZeroEx(web3.currentProvider);
    Exchange.at(order.exchangeContractAddress).ZRX_TOKEN_CONTRACT.call().then((zrx_address) => {
        return Promise.all([
            zeroEx.token.transferAsync(order.takerTokenAddress, web3.eth.accounts[0], web3.eth.accounts[1], order.takerTokenAmount.mul(1.1)),
            zeroEx.token.transferAsync(zrx_address, web3.eth.accounts[0], web3.eth.accounts[1], order.takerFee.mul(1.1)),
            zeroEx.token.setProxyAllowanceAsync(zrx_address, web3.eth.accounts[1], order.takerFee.mul(1.1)),
            zeroEx.token.setProxyAllowanceAsync(order.takerTokenAddress, web3.eth.accounts[1], order.takerTokenAmount.mul(1.1)),
            zeroEx.token.setProxyAllowanceAsync(zrx_address, web3.eth.accounts[0], order.makerFee.mul(1.1)),
            zeroEx.token.setProxyAllowanceAsync(order.makerTokenAddress, web3.eth.accounts[0], order.makerTokenAmount.mul(1.1)),
        ]);
    }).then(() => {
        return zeroEx.exchange.fillOrderAsync(order, order.takerTokenAmount.div(2), false, web3.eth.accounts[1]);
    }).then(done).catch((error) => {
        console.log(error);
    });

}

function fixBn(order) {
    order.makerTokenAmount = new bignumber(order.makerTokenAmount);
    order.takerTokenAmount = new bignumber(order.takerTokenAmount);
    order.makerFee = new bignumber(order.makerFee);
    order.takerFee = new bignumber(order.takerFee);
    order.expirationUnixTimestampSec = new bignumber(order.expirationUnixTimestampSec);
    return order
}
