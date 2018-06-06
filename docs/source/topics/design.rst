Technical Design
================

OpenRelay was designed to be a scalable relayer for the 0x protocol. It uses a
microservices architecture and is written entirely in Go. While OpenRelay
consists of around 20 microservices, they are all relatively light weight and
can run comfortably on a modern laptop.

This document assumes that you are familiar with the basic architecture of the
0x protocol.

High Level Architecture
-----------------------

OpenRelay was designed to use a message broker to pass messages between
microservices. This allows us to scale the individual services as needed, and
allows for services to get backed up on work instead of rejecting subsequent
requests (and once the services scale up, they can collectively handle the
workload to help catch up). At present, OpenRelay.xyz is deployed on Amazon Web
Services, but should be able to run on any infrastructure provided it has
access to a Redis server and either a MySQL or PostgreSQL server.

Language Choice
...............

When OpenRelay was originally built, most of the services were written in Go,
while a few were written in Python and JavaScript. Since that time, we have
replace all of the Python and JavaScript services with versions written in Go.

Go was chosen for a number of reasons. When processing a high volume of orders,
OpenRelay is a rather computationally intensive application. Where Python and
JavaScript are interpreted languages, Go runs as a native application, giving
it a performance boost that helps keep computational requirements down.


Additionally, OpenRelay is designed to run within Docker, and while Python and
JavaScript both require minimal Linux distributions in the Docker container, Go
binaries can function as the only file in the container. OpenRelay's containers
have only a single, statically compiled binary. They have no shell, and no
executables other than the Microservices they're built to run. This helps keep
a minimal attack surface, as even if an outside attacker managed to exploit one
of our services, the container gives them no other commands to execute to do
things like get a shell.


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

Many services are capable of supporting multiple concurrent messaging routes.
That is, a single service instance can consume messages from `queue://ingest`
and publish those messages to `queue://fundcheck` while also pulling messages
from `queue://released` and publishing those messages to `queue://recheck`.
This reduces the number of instances required for each service.

Ingest Service
^^^^^^^^^^^^^^

The first service in the pipeline is the ingest service. The ingest service
listens for HTTP requests, parses out provided Orders, performs basic
validation on those orders, and places them on a message queue to be further
validated and eventually indexed.

The ingest service provides the POST `/v0/order` and `/v0/fees` APIs from the
`Standard Relayer API <https://github.com/0xProject/standard-relayer-api/blob/master/http/v0.md>`_.

* **Classification**: External

Fill Updater
^^^^^^^^^^^^

The fill updater receives orders from a message queue and communicates with an
Ethereum node to determine whether the order has been fully or partially
filled. It updates the order, and publishes the order on an outbound channel.

The fill updater also watches a topic provided by the Fill Monitor, and uses it
to maintain a bloom filter of orders that have been filled on the Exchange
contract. The bloom filter allows the fill updater to check whether orders are
likely to have been filled without making any API calls on a per-order basis
unless the order appears in the bloom filter.

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
cancelled, based on the values set in the Fill Updater service. If a given
maker submits multiple orders for the same maker token in the span of a single
block, the fund validator will only check balances and allowances once per
block.

* **Classification**: Internal

Simple Relay
^^^^^^^^^^^^

The simple relay takes messages from one topic or queue, and places them onto
one or more other topics or queues. This is a very basic function, and is used
in a few places throughout the OpenRelay pipeline.

* **Classification**: Internal

Delay relays
^^^^^^^^^^^^

On OpenRelay.xyz, premium members will get immediate access to orders, while
unpaid users get access to orders two Ethereum blocks after the order is
received. This is controlled by delay relays. Delay relays have three
parameters:

* Input Queue: The queue where messages wait to be relayed
* Output Queue: The queue or topic where messages will be relayed
* Signal Queue: The queue or topic where messages are received to signal the
  release of messages from the input queue to the output queue.

  * **Classification**: Internal

SQL Order Index
^^^^^^^^^^^^^^^^

The SQL Order Index provides more robust search capabilities, and is managed by
several indexing. OpenRelay supports both MySQL and PostgreSQL databases for
the SQL order index.

SQL Search API
^^^^^^^^^^^^^^^

The SQL search API allows internet based users to query the Order Index for
orders.

The SQL search API provides the GET `/v0/token_pairs`, `/v0/orders`,
`/v0/order/${order_hash}`, and `/v0/orderbook` endpoints from the
`Standard Relayer API <https://github.com/0xProject/standard-relayer-api/blob/master/http/v0.md>`_.

* **Classification**: External

Block Monitor
^^^^^^^^^^^^^

The block monitor polls the blockchain watching for new blocks. When it finds
them, it emits a message including:

* The block number
* The block hash
* The block bloom filter

Other services that need to monitor for on-chain events can consume the output
of the block monitor, and use the bloom filter to check for events of interest
to that monitoring service before having to make RPC calls to query for events.

* **Classification**: Internal


Exchange Monitor
^^^^^^^^^^^^^^^^

The exchange monitor service consumes messages from the block monitor service,
watching for LogFill and LogCancel events on the Exchange contract. The
exchange monitor emits events containing the order hash, and either the
`filledTakerTokenAmount` or the `cancelledTakerTokenAmount`, depending on the
event type.

* **Classification**: Internal

Spend Monitor
^^^^^^^^^^^^^

The spend monitor service consumes messages from the block monitor service,
watching for ERC20 Transfer events. When it finds a transfer event, it queries
for both the allowance a user has set for the 0x Token Transfer Proxy, and the
current balance of that user for that token. The lesser of the allowance and
the balance will determine whether this event has changed whether or not an
order by the spender is still fillable, so the Spend Monitor emits a Spend
Record message.

* **Classification**: Internal

Allowance Monitor
^^^^^^^^^^^^^^^^^

The allowance monitor service consumes messages from the block monitor service,
watching for ERC20 Approve events where the 0x Token Transfer Proxy is the
approved spender. When it finds an approve event, it emits a Spend Record
message indicating the current allowance level.

* **Classification**: Internal

Spend Indexer
^^^^^^^^^^^^^

The spend indexer consumes Spend Record messages emitted by the Spend Monitor
and the Allowance monitor. When it receives messages, it sends a query to the
database to mark unfillable any orders where the maker matches the spender on
the spend record and where the remining maker fill amount exceeds the spender's
remaining balance.

* **Classification**: Internal

SQL Fill Indexer
^^^^^^^^^^^^^^^^

The fill indexer consumes the messages emitted by the Exchange Monitor and uses
them to update the SQL Order Index with the cancelled and filled amounts of
each received message. After updating the record, if the cancelled + filled
amounts equal the total taker amount, the record is marked as filled, which
removes it from most of the search endpoints.

* **Classification**: Internal

Exchange Splitter
^^^^^^^^^^^^^^^^^

In order to support multiple networks (Mainnet and the Ropsten testnet), we
have deployed multiple instances of some of these services, some to manage
mainnet and some to manage testnet. The ingest service is shared by both
deployments. The Exchange splitter picks up orders from the ingest service, and
routes them to the appropriate pipeline for verification. After verification,
both pipelines feed their results back to the same indexer service.

* **Classification**: Internal

Ethereum Nodes
..............

At present OpenRelay uses Infura for Ethereum state lookups. We are working
hard on a solution to bring Ethereum node hosting in-house, but the operational
complexities have proven challenging, and we did not want to hold up use of
OpenRelay while work through these challenges.
