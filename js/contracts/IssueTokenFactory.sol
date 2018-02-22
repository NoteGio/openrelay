pragma solidity ^0.4.11;

import "./base/IssueToken.sol";

contract IssueTokenFactory {

    function IssueTokenFactory() {
    }

    function newToken(string _name, string _symbol) public returns (address newToken) {
        IssueToken c = (new IssueToken(_name, _symbol));
        c.transferOwnership(msg.sender);
        return address(c);
    }
}
