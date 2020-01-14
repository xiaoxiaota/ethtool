// pragma solidity ^0.4.25;
pragma solidity ^0.5.11;

contract Storage {
  uint256 storedData;
  function set(uint256 data) public {
    storedData = data;
  }
  function Get() public view returns (uint256) {
    return storedData;
  }
}
