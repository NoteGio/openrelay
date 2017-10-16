Usage Guide
===========


0x Standard Relayer API
-----------------------

OpenRelay implements a draft of the
`0x Standard Relayer API <https://github.com/0xProject/standard-relayer-api/blob/e6a962e5d6d2f8c0ca53c3c6e6aa6ba34bc451bc/README.md>`_,
with a few exceptions.

* The 0x Standard Relayer API is still a draft, and has changed since OpenRelay
  made our initial implementation. Since we need something resembling a stable
  API now, but hope to correctly implement the standard once it is finalized,
  we are namespacing our api with `/v0.0/` instead of `/v0/`, and we will
  increment the decimal as we get closer to the standard specification. Once
  the standard is finalized, we will provide a `/v0/` endpoint.
* OpenRelay does not implement the `/v0.0/token_pairs` Endpoint. OpenRelay will
  support any valid ERC-20 token, in any combination of pairs. OpenRelay is not
  concerned with token symbols (using only the contract addresses for
  referencing tokens), nor is it concerned with decimals. The token pairs
  endpoint may be implemented in the future, but is not on the immediate
  roadmap.
* OpenRelay does not implement several of the parameters for the `/v0.0/orders`
  endpoint. OpenRelay's current storage engine is somewhat limited in terms of
  allowing complex searches. Replacing the storage engine is on the medium-term
  roadmap, which should enable it to support all of the parameters for the
  `/v0.0/orders`. As of now, the supported parameters are:

   * **makerTokenAddress**
   * **takerTokenAddress**
   * **limit**

  When both makerTokenAddress and takerTokenAddress are specified, the results
  are sorted by price (takerAmount / makerAmount). Otherwise the results are
  unsorted.
* The 0x Standard Relayer API only allows for `/v0.0/fees` to allow a single
  configuration of valid fees. OpenRelay requires a set total fee, but does not
  dictate the maker's share or the taker's share. The response of `/v0.0/fees`
  will specify the total fee as the maker's share, but `/v0.0/order` will accept
  any distribution such that makerFee and takerFee total the specified value.
* The `/v0.0/order` endpoint will return a 202 status code upon success instead
  of a 201 status code. OpenRelay queues orders for additional validation and
  processing between receipt and listing it in the order book, thus it is more
  correct to respond that the order was "Accepted" rather than "Created".
* When making GET requests against `/v0.0/orders` or `/v0.0/order/[order_hash]`,
  a blockhash parameter is required. If none is provided, you will get a 307
  redirect to the same resource with the blockhash parameter provided. This is
  included to improve caching, and leverages the fact that the orderbook will
  only change as Ethereum blocks complete.

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
OpenRelay and included with the response from `/v0.0/orders`.

To submit a 377 byte order using the binary format, POST to `/v0.0/order` with
the header `Content-Type: application/octet-stream`.

To retrieve a list of orders in the binary format, GET `/v0.0/orders` with the
header `Accept: application/octet-stream`. Othere parameters remain the same.
Note that the response stream may consist of numerous orders; every 441 bytes
marks the start of a new order.
