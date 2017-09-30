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

Python is used for interacting with DynamoDB, which is the current persistent
store for order information. The OpenRelay team was already familiar and
comfortable with Python's DynamoDB libraries, so it seemed an easy choice to
get started. As we have gotten further into the project, we have deemed that
DynamoDB will need to be replaced to achieve some of our additional goals, and
replacing DynamoDB will likely involve eliminating Python from the OpenRelay
stack.

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

* **Language**: Go
* **Classification**: External

Fund Validation
^^^^^^^^^^^^^^^

The fund validation receives orders from a message queue, and communicates with
an Ethereum node to validate that:

* The maker of the order has adequate funds to complete the order (both the
  maker token and the fee token).
* The maker of the order has set allowances for the token transfer proxy to
  complete the order.

It also checks to see if the order has been all or partially filled or
cancelled, and includes an accurate TakerTokenAmountFilled as it relays the
message.

* **Language**: Go
* **Classification**: Internal

Simple Relay
^^^^^^^^^^^^

The simple relay takes messages from one topic or queue, and places them onto
another topic or queue. This is a very basic function, and is used in a few
places throughout the OpenRelay pipeline.

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

Indexer
^^^^^^^

After messages make it through the delay relays, they are picked up by the
indexer. The indexer parses the messages, and stores them in the order index
for subsequent retrieval.

* **Language**: Python
* **Classification**: Internal

Order Index
^^^^^^^^^^^

The order index is responsible for persistent storage of orders that have been
validated. The order index is currently implemented with DynamoDB. DynamoDB is
rather limiting, as it offers a limited number of indexable fields, and lacks
query planning or index intersection. The Order Index will likely be replaced
in an upcoming release.

Search API
^^^^^^^^^^

The search API allows internet-based users to query the Order Index for orders.
At present, it allows searching for orders by:

* Maker Token
* Taker Token
* Token Pair

* **Language**: Python
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

Fill Indexer
^^^^^^^^^^^^

The fill indexer consumes the messages emitted by the Exchange Monitor and uses
them to update the Order Index with the cancelled and filled amounts of each
received message. After updating the record, if the cancelled + filled amounts
equal the total taker amount, the record is deleted.

* **Language**: Python
* **Classification**: Internal

Ethereum Nodes
..............

OpenRelay uses Ethereum clients for two different purposes, and uses separate
containers for each purpose to help manage the load.

Monitoring Node
^^^^^^^^^^^^^^^

OpenRelay has exactly one monitoring node at any given time. At present, the
monitoring node is a light geth client. Filters can be installed for each event
that needs to be processed, and the monitoring node will provide a feed of
events.

Because filters need to be installed on the monitoring node, monitoring nodes
cannot be load balanced. Load is not a huge concern for monitoring nodes, as
the load on the monitoring node is relatively constant, and doesn't scale much
with increased relay activity.

For failover, the monitoring node is fronted with HAProxy, which will fail-over
to the Standby Node if the monitoring node becomes unavailable. This will cause
filters to be dropped, so monitoring services need to be prepared to re-install
their filters.

State Nodes
^^^^^^^^^^^

State nodes are used for looking up token balances, allowances, order fill
amounts, and order cancellations. There can be an arbitrary number of state
nodes and they can be load balanced, as each state request is independent. The
state node is also fronted by HAProxy, and will fall back to the standby node
if needed.

Standby Node
^^^^^^^^^^^^

A single standby node is the fallback for both the monitoring node and the
standby node. Under normal circumstances it should see little internal traffic.
It exists to handle traffic if the monitoring node or all of the state nodes
die, until the monitoring node can be restored.
