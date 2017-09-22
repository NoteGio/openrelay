module.exports = {
  networks: {
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
    },
    main: {
      host: "monitor-haproxy",
      port: 8545,
      network_id: "*"
    }
  }
};
