pragma solidity ^0.4.11;

import "./StandardToken.sol";
import "./Ownable.sol";

contract IssueToken is StandardToken, Ownable {
    function issue(uint256 quantity) public onlyOwner {
        require(totalSupply + quantity > totalSupply);
        require(balances[owner] + quantity > balances[owner]);
        balances[owner] += quantity;
        totalSupply += quantity;
    }
}
