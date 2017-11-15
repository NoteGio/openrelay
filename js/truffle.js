var ProviderEngine = require("web3-provider-engine");
var WalletSubprovider = require('web3-provider-engine/subproviders/wallet.js');
var Web3Subprovider = require("web3-provider-engine/subproviders/web3.js");
var Web3 = require("web3");
const FilterSubprovider = require('web3-provider-engine/subproviders/filters.js');

function getEngine() {
    if(!process.env.USE_FILTER_PROVIDER) {
        return new Web3.providers.HttpProvider(process.env.ETHEREUM_URL)
    }
    var engine = new ProviderEngine();
    engine.addProvider(new FilterSubprovider());
    engine.addProvider(new Web3Subprovider(new Web3.providers.HttpProvider(process.env.ETHEREUM_URL)));
    engine.start();
    return engine
}

module.exports = {
  networks: {
    main: {
      provider: getEngine,
      network_id: "*"
    },
    development: {
      host: "localhost",
      port: 8546,
      network_id: "*" // Match any network id
    },
    testnet: {
      host: "ethnode",
      port: 8545,
      network_id: "*"
    },
    parity: {
      host: "172.17.0.4",
      port: 8545,
      network_id: "*" // Match any network id
    }
  }
};
