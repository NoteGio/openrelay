Technical Design
================

OpenRelay was designed to be a scalable relayer for the 0x protocol. It uses a
microservices architecture, with components written in Go, JavaScript and
Python. While OpenRelay consists of around 20 microservices, they are all
relatively light weight and can run comfortably on a modern laptop.

This document assumes that you are familiar with the basic architecture of the
0x protocl.

High Level Architecture
-----------------------

OpenRelay was designed to use a message broker to pass messages between
microservices. This allows us to scale the individual services as needed, and
allows for services to get backed up on work instead of rejecting subsequent
requests (and once the services scale up, they can collectively handle the
workload to help catch up). At present, OpenRelay is designed to be deployed on
Amazon Web Services.

Language Choices
................

Most of the services within OpenRelay are written in Go. Go was chosen
primarily because of its speed and efficiency. Go is the default choice for
microservices, with JavaScript and Python each filling a specific niche.

JavaScript is used for monitoring events on the Ethereum blockchain. Go lacks
advanced libraries for monitoring for contract events, and JavaScript's Web3
libraries provide a nice interface for monitoring events. As the JavaScript
microservices find events via Web3, the events are placed on message queues for
subsequent processing by other services.

Python is used for interacting with DynamoDB, which is the legacy persistent
store for order information. It has been replaced with a PostgreSQL order book,
which is managed by Go services, using GORM to interact with the database
layer.

Service Classifications
.......................

Microservices are classified as either internal or externally facing.

External services handle requests from users on the Internet. This impacts the
scaling considerations (externally facing services could see sharp spikes in
usage) and security considerations (externally facing services are more likely
to get malicious inputs from third parties). External services have read-only
access to the Order index, and can only get messages into the index by passing
messages to internal services for validation.

Internal services either monitor for Ethereum events, or pass messages on the
message broker. They can generally assume that the messages they receive will
be well formed, and can process messages from the message broker as best as
they are able without signficantly impacting the end-user. Internal services
may have read/write access to the order index.

Services
........

Message Broker
^^^^^^^^^^^^^^

The current message broker used by OpenRelay is Redis. The message broker is
the backbone of OpenRelay, and Redis was chosen because it is reasonably
scalable, and hosted Redis is available with Amazon Elasticache. Redis may
eventually be replaced by Kafka or another message broker, but Redis seemed
like a good option to get OpenRelay started with minimal fuss.

Ingest Service
^^^^^^^^^^^^^^

The first service in the pipeline is the ingest service. The ingest service
listens for HTTP requests, parses out provided Orders, performs basic
validation on those orders, and places them on a message queue to be further
validated and eventually indexed.

The ingest service provides the POST `/v0/order` and `/v0/fees` APIs from the
[Standard Relayer API](https://github.com/0xProject/standard-relayer-api/blob/master/http/v0.md).

* **Language**: Go
* **Classification**: External

Fill Updater
^^^^^^^^^^^^

The fill updater receives orders from a message queue and communicates with an
Ethereum node to determine whether the order has been fully or partially
filled. It updates the order, and publishes the order on an outbound channel.

* **Language**: Go
* **Classification**: Internal

Fund Validation
^^^^^^^^^^^^^^^

The fund validation receives orders from a message queue, and communicates with
an Ethereum node to validate that:

* The maker of the order has adequate funds to complete the order (both the
  maker token and the fee token).
* The maker of the order has set allowances for the token transfer proxy to
  complete the order.

It also takes into consideration whether the order has been partially filled or
cancelled, based on the values set in the Fill Updater service.

* **Language**: Go
* **Classification**: Internal

Simple Relay
^^^^^^^^^^^^

The simple relay takes messages from one topic or queue, and places them onto
one or more other topics or queues. This is a very basic function, and is used
in a few places throughout the OpenRelay pipeline.

* **Language**: Go
* **Classification**: Internal

Delay relays
^^^^^^^^^^^^

On OpenRelay.xyz, premium members get immediate access to orders, while unpaid
users get access to orders two Ethereum blocks after the order is received.
This is controlled by delay relays. Delay relays have three parameters:

* Input Queue: The queue where messages wait to be relayed
* Output Queue: The queue or topic where messages will be relayed
* Signal Queue: The queue or topic where messages are received to signal the
  release of messages from the input queue to the output queue.

  * **Language**: Go
  * **Classification**: Internal

DynamoDB Indexer
^^^^^^^^^^^^^^^^

After messages make it through the delay relays, they are picked up by the
indexer. The indexer parses the messages, and stores them in the order index
for subsequent retrieval.

* **Language**: Python
* **Classification**: Internal
* **Deprecated**: The DynamoDB index is still supported, but is not receiving
  many updates.

DynamoDB Order Index
^^^^^^^^^^^^^^^^^^^^

The order index is responsible for persistent storage of orders that have been
validated. The order index is currently implemented with DynamoDB. DynamoDB is
rather limiting, as it offers a limited number of indexable fields, and lacks
query planning or index intersection, which is why it is being replaced with a
PostgreSQL order index.

* **Deprecated**: The DynamoDB index is still supported, but is not receiving
  many pudates.

PostgreSQL Order Index
^^^^^^^^^^^^^^^^^^^^^^

The PostgreSQL order index currently lives alongside the DynamoDB index. It
provides more robust search capabilities, and is managed by Go microservices
instead of Python microservices.

DyanmoDB Search API
^^^^^^^^^^^^^^^^^^^

The search API allows internet-based users to query the Order Index for orders.
At present, it allows searching for orders by:

* Maker Token
* Taker Token
* Token Pair

* **Language**: Python
* **Classification**: External

PostgreSQL Search API
^^^^^^^^^^^^^^^^^^^^^

The PostgreSQL search API allows internet based users to query the Order Index
for orders.

The PostgreSQL search API provides the GET `/v0/token_pairs`, `/v0/orders`,
`/v0/order/${order_hash}`, and `/v0/orderbook` endpoints from the [Standard Relayer API](https://github.com/0xProject/standard-relayer-api/blob/master/http/v0.md).

* **Language**: Go
* **Classification**: External

Block Monitor
^^^^^^^^^^^^^

The block monitor service emits events every time a block is mined in the
Ethereum blockchain. Events can be sent either to topics or queues.

* **Language**: JavaScript
* **Classification**: Internal


Exchange Monitor
^^^^^^^^^^^^^^^^

The exchange monitor service listens for LogFill and LogCancel events on the
Exchange contract, filtering for events where the feeRecipient is a an
authorized fee recipient for the relay. The exchange monitor emits events
containing the order hash, and either the `filledTakerTokenAmount` or the
`cancelledTakerTokenAmount`, depending on the event type.

* **Language**: JavaScript
* **Classification**: Internal

DynamoDB Fill Indexer
^^^^^^^^^^^^^^^^^^^^^

The fill indexer consumes the messages emitted by the Exchange Monitor and uses
them to update the DyanmoDB Order Index with the cancelled and filled amounts
of each received message. After updating the record, if the cancelled + filled
amounts equal the total taker amount, the record is deleted.

* **Language**: Python
* **Classification**: Internal

PostgreSQL Fill Indexer
^^^^^^^^^^^^^^^^^^^^^^^

The fill indexer consumes the messages emitted by the Exchange Monitor and uses
them to update the PostgreSQL Order Index with the cancelled and filled amounts
of each received message. After updating the record, if the cancelled + filled
amounts equal the total taker amount, the record is marked as filled, which
removes it from most of the search endpoints.

* **Language**: Go
* **Classification**: Internal

Exchange Splitter
^^^^^^^^^^^^^^^^^

In order to support multiple networks (Mainnet and the Ropsten testnet), we
have deployed multiple instances of most of these services, some to manage
mainnet and some to manage testnet. The ingest service is shared by both
deployments. The Exchange splitter picks up orders from the ingest service, and
routes them to the appropriate pipeline for verification. After verification,
both pipelines feed their results back to the same indexer service.

* **Language**: Go
* **Classification**: Internal

Ethereum Nodes
..............

At present OpenRelay uses Infura for Ethereum state lookups. We are working
hard on a solution to bring Ethereum node hosting in-house, but the operational
complexities have proven challenging, and we did not want to hold up use of
OpenRelay while work through these challenges.
