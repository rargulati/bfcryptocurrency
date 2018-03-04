pragma solidity ^0.4.19;

contract GamblingGame {

    struct Player {
        address payout;
        bytes32 commit;
        uint reveal;
    }
    
    uint8 playersRevealed;
    Player[] players;

    function EnterGame(bytes32 _commit) public payable{
        require(msg.value == 1 ether);
        
        address _player = msg.sender;

        _createPlayer(_player, _commit);
    }
        
    function _createPlayer(address _player, bytes32 _commit) private {
        players.push(Player(_player, _commit, 0));
    } 
    
    function RevealCommit(uint _index, uint _nonce, uint _number) private {
        require(players.length == 2);
        address _player = msg.sender;

        require(uint(_nonce) != 0);
        require(players[_index].payout == _player);
        
        require(keccak256(_index, _nonce, _number) == players[_index].commit);
        players[_index].reveal = _number;
        
        playersRevealed++;
        
        if (playersRevealed == 2) {
            uint sum = players[0].reveal + players[1].reveal;
            uint winner = sum % 2;
            players[winner].payout.transfer(2 ether);
        }
    }

    function Keccak(uint _index, uint _nonce, uint _number) public returns (bytes32) {
      return keccak256(_index, _nonce, _number);
    }
}
