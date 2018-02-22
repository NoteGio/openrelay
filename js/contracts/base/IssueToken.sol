pragma solidity ^0.4.11;

import "./StandardToken.sol";
import "./Ownable.sol";

contract IssueToken is StandardToken, Ownable {
    string public name;
    string public symbol;
    uint8 public constant decimals = 18;
    function IssueToken(string _name, string _symbol) Ownable() {
        name = _name;
        symbol = _symbol;
    }

    function issue(uint256 quantity) public onlyOwner {
        require(totalSupply + quantity > totalSupply);
        require(balances[owner] + quantity > balances[owner]);
        balances[owner] += quantity;
        totalSupply += quantity;
    }
}
