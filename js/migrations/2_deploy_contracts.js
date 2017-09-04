var IssueTokenFactory = artifacts.require("./IssueTokenFactory.sol");

module.exports = function(deployer) {
  deployer.deploy(IssueTokenFactory);
};
