Usage Guide
===========


0x Standard Relayer API
-----------------------

OpenRelay implements the
`0x Standard Relayer API <https://github.com/0xProject/standard-relayer-api/blob/master/http/v0.md>`_.

You can interface with it using `OpenRelay.js<https://github.com/NoteGio/openrelay.js>`_,
`0x Connect<https://0xproject.com/docs/connect>`_, or any other client for the 0x standard relayer API.


Binary Format
-------------

In addition to the JSON format from the 0x Standard Relayer API, OpenRelay
supports a binary format which uses less bandwidth and has shorter processing
times. That binary format consists of:

    ============ ========== ========================== =========
     Start Byte   End Byte           Content             Format
    ============ ========== ========================== =========
    0            20         exchangeContractAddress    [20]byte
    20           40         maker                      [20]byte
    40           60         taker                      [20]byte
    60           80         makerTokenAddress          [20]byte
    80           100        takerTokenAddress          [20]byte
    100          120        feeRecipient               [20]byte
    120          152        makerTokenAmount           uint256
    152          184        takerTokenAmount           uint256
    184          216        makerFee                   uint256
    216          248        takerFee                   uint256
    248          280        expirationUnixTimestampSec uint256
    280          312        salt                       uint256
    312                     ecSignature.v              byte
    313          345        ecSignature.r              [32]byte
    345          377        ecSignature.s              [32]byte
    377          409        takerTokenAmountFilled     uint256
    409          441        takerTokenAmountCancelled  uint256
    ============ ========== ========================== =========

Bytes 0-312 can be hashed with keccak256, and the signature in bytes 313-377
should match the hash.

Only the first 377 bytes should be submitted with an order. The
takerTokenAmountFilled and takerTokenAmountCancelled will be populated by
OpenRelay and included with the response from `/v0/orders` if you request the
binary format.

To submit a 377 byte order using the binary format, POST to `/v0/order` with
the header `Content-Type: application/octet-stream`.

To retrieve a list of orders in the binary format, GET `/v0.0/orders` with the
header `Accept: application/octet-stream`. Othere parameters remain the same.
Note that the response stream may consist of numerous orders; every 441 bytes
marks the start of a new order.
