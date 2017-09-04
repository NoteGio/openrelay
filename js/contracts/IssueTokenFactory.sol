pragma solidity ^0.4.11;

import "./base/IssueToken.sol";

contract IssueTokenFactory {

    function IssueTokenFactory() {
    }

    function newToken() public returns (address newToken) {
        IssueToken c = (new IssueToken());
        c.transferOwnership(msg.sender);
        return address(c);
    }
}
