Tutorials
=========

The following tutorials are intended to help you get started interacting with
OpenRelay and the :term:`0x protocol` using JavaScript. Any of these things
could be done in other languages as well, but the best available tooling is
written in JavaScript, so that's where we're starting out. We would welcome
pull requests to add documentation for other languages!

Prerequisites
-------------

Getting Started with Web3
.........................

To interact with the Ethereum Blockchain, you'll need a :term:`Web3` provider.
If you are building a DApp for other users to interact with, we'd recommend
starting off with `Metamask <https://metamask.io/>`_. If you are building a
backend system, take a look at the
`Node JS Package <https://www.npmjs.com/package/web3>`_, or
`Truffle <http://truffleframework.com/docs/>`_. Truffle is largely geared
towards contract development, but generally provides a good toolchain for
interacting with existing contracts.

The tutorials that follow will assume that you have a web3 object defined and
ready to interact with the blockchain.

0x.js
.....

We will also make heavy use of the `0x.js <https://0xproject.com/docs/0xjs>`_
library. If you are using Node.js, install it with

.. code-block:: bash

    npm install 0x.js --save

If you are including it in an HTML application, include it with

.. code-block:: html

    <script type="text/javascript" src="0x.js"></script>

BigNumber.js
............

Finally, you will need the BigNumber library, as JavaScript's numbers can't
accurately represent numbers at the scale often used for an :term:`ERC20 token`.

If you're using NPM it was installed when you insalled 0x.js. You can import it
with

.. code-block:: javascript

    var BigNumber = require('bignumber.js');


If you are
creating including it in an HTML application, add

.. code-block:: html

    <script src='relative/path/to/bignumber.js'></script>


Issue an Order on OpenRelay.xyz
-------------------------------

Creating The Order
..................

An order starts out life as a simple JavaScript object. For the purpose of
this tutorial, we will make an :term:`offer` of 1000 ZRX for 0.7 :term:`WETH`.
We will assume that the maker of the order is the default account on our Web3
object, and that taker is `null`, allowing anyone to fill the order.

The 0x wiki has a `list of deployed contract addresses <https://0xproject.com/wiki#Deployed-Addresses>`_,
which we will use below.

We'll start constructing our order object:

.. code-block:: javascript

    // Order will be valid for 24 hours
    var duration = 60*60*24;

    var order = {
        // The default web3 account address
        maker: web3.eth.accounts[0],
        // Anyone may fill the order
        taker: "0x0000000000000000000000000000000000000000",
        // The ZRX token contract on mainnet
        makerTokenAddress: "0xe41d2489571d322189246dafa5ebde1f4699f498",
        // The WETH token contract on mainnet
        takerTokenAddress: "0x2956356cd2a2bf3202f771f50d3d14a367b48070",
        // A BigNumber objecct of 1000 ZRX. The base unit of ZRX has 18
        // decimal places, the number here is 10^18 bigger than the
        // base unit.
        makerTokenAmount: new BigNumber("1000000000000000000000"),
        // A BigNumber objecct of 0.7 WETH. The base unit of WETH has
        // 18 decimal places, the number here is 10^18 bigger than the
        // base unit.
        takerTokenAmount: new BigNumber("700000000000000000"),
        // Add the duration (above) to the current time to get the unix
        // timestamp
        expirationUnixTimestampSec: parseInt(
                (new Date().getTime()/1000) + duration
            ).toString(),
        // We need a random salt to distinguish different orders made by
        // the same user for the same quantities of the same tokens.
        salt: ZeroEx.ZeroEx.generatePseudoRandomSalt()
    }

At this point, we have most of our order object, but it isn't complete. We need
to look up the exchange contract address. At the time of this writing, this
will be `0x12459c951127e0c374ff9105dda097662a027093` on mainnet, but by looking
it up with the ZeroEx library we are future-proofing against exchange contract
upgrades.

.. code-block:: javascript

    var zeroEx = new ZeroEx.ZeroEx(web3.currentProvider);
    var addressPromise = zeroEx.exchange.getContractAddressAsync().then(
        (address) => {
            order.exchangeContractAddress = address;
        }
    )

Then we need to determine what fees the relayer will require. This example uses
the `request promise library <https://github.com/request/request-promise>`_ for
compatibility between NodeJS and browser implementations. You can use your
choice of HTTP client libraries. For this example we'll get fees from the
public OpenRelay API, but this example should work with a private relay, or
any other relay implementing the
`standard relayer API specification <https://github.com/0xProject/standard-relayer-api>`_.

.. code-block:: javascript

    //
    var openrelayBaseURL = "https://api.openrelay.xyz";
    var feePromise = rp({
        method: 'POST',
        uri: openrelayBaseURL + "/v0/fees",
        body: order,
        json: true,
    }).then((feeResponse) => {
        // Convert the makerFee and takerFee into BigNumbers
        order.makerFee = new BigNumber(feeResponse.makerFee);
        order.takerFee = new BigNumber(feeResponse.takerFee);
        // The fee API tells us what taker to specify
        order.taker = feeResponse.takerToSpecify;
        // The fee API tells us what feeRecipient to specify
        order.feeRecipient = feeResponse.feeRecipient;
    })

It's worth noting that OpenRelay.xyz will accept any distribution of fees
between the maker and the taker so long as the total fee meets the minimum, but
most other relayers require you to stick with the fees they specify.

Signing The Order
.................

Once the order is defined, we need to sign it with our Ethereum account. This
is what the :term:`Exchange contract` needs as proof that the maker intended
to authorize an order.

.. code-block:: javascript

    Promise.all([addressPromise, feePromise]).then(() => {
        // Once those promises have resolved, our order is ready to be signed
        var orderHash = ZeroEx.ZeroEx.getOrderHashHex(order);
        return zeroEx.signOrderHashAsync(orderHash, order.maker);
    }).then((signature) => {
      order.ecSignature = signature;
      return order;
    });

Note that if you are using Metamask or Mist, your users will be prompted to
approve signing the message. This is a security measure to ensure that web
applications can't maliciously sign messages on behalf of their users.

Approving The Transfer
......................

Now that we have a signed order object, we're almost ready to submit it to a
relay. But the relay will reject our order if we don't have the funds required
to fill the order. We need to allow the :term:`Token Proxy` to transfer funds
on our behalf, which it will do when someone tries to fill our order. There are
two ways to handle this.

Set Exact Allowances
^^^^^^^^^^^^^^^^^^^^
Option one is to approve exactly the amount of funds required for this order.
This is a low risk option, and even if the Token Transfer Proxy contract is
somehow compromised in the future, it will only be able to transfer the funds
you intended to transfer as part of an order. To do this:

.. code-block:: javascript

    var makerAllowance = zeroEx.token.setProxyAllowanceAsync(
        order.makerToken,
        order.maker,
        order.makerTokenAmount
    );

You also need to set allowances for the fees:

.. code-block:: javascript

    var feeAllowance = zeroEx.token.setProxyAllowanceAsync(
        "0xe41d2489571d322189246dafa5ebde1f4699f498",
        order.maker,
        order.makerFee
    );

Note that each time you set the proxy allowance, you are setting the exact
value; it does not add allowances together, so if you have multiple orders at
the same time you need to track all of them and set allowances cumulatively.

Set Unlimited Allowances
^^^^^^^^^^^^^^^^^^^^^^^^

The second option is to allow the Token Proxy to transfer unlimited quantities
of a given token on your behalf. This isn't a very high risk proposition, as
the Token Proxy is an
`open source Ethereum contract <https://etherscan.io/address/0x12459c951127e0c374ff9105dda097662a027093#code>`,
and its code can't be changed without changing the address. However if you're
concerned that it may somehow be compromised in the future, it's somewhat safer
to authorize smaller quantities at a time.

The benefits of unlimited allowance are:

* You don't have to keep track of open orders to make sure your allowances add
  up to all open orders.
* You don't have to pay for gas every time you want to increment your
  allowances.
* In the future, you won't have to wait for allowance transactions to complete
  before you can fill an order.

To set unlimited allowances for your makerToken and the ZRX Token, run:

.. code-block:: javascript

    var makerAllowance = zeroEx.token.setUnlimitedProxyAllowanceAsync(
        order.makerToken,
        order.maker,
    );
    var feeAllowance = zeroEx.token.setUnlimitedProxyAllowanceAsync(
        "0xe41d2489571d322189246dafa5ebde1f4699f498",
        order.maker,
    );

Submit The Order to OpenRelay.xyz
.................................

Once your allowance transactions have been included in a block, you're ready to
submit your order to the OpenRelay API. Once again using the request-promise
library:

.. code-block:: javascript

    rp({
        method: 'POST',
        uri: openrelayBaseURL + "/v0/order",
        body: order,
        json: true,
    })

If you've done everything right up to this point, this should return a 202
'Accepted' status code. OpenRelay will then do some additional validations,
double checking that you have the funds necessary to :term:`fill` the order.
Your order should be posted to the order book after two ethereum blocks have
been mined.

Find Orders on OpenRelay
------------------------

Searching the Order Book
........................

Searching the OpenRelay :term:`order book` is fairly simple. Using the
request-promise library from the previous tutorial,

.. code-block:: javascript

    rp({
        method: 'GET',
        uri: openrelayBaseURL + "/v0/orders",
        json: true,
    }).then((orders) => {
        for(var order of orders) {
            // Manipulate order object
        }
    })

The above search will return any :term:`order` in the order book. If you want
to narrow the search, you can use the following search parameters:

* makerTokenAddress - Filter for orders where makerTokenAddress matches the
  specified value.
* takerTokenAddress - Filter for orders where takerTokenAddress matches the
  specified value.

You may specify zero, one, or both of the makerTokenAddress and
takerTokenAddress parameters. To find our order from the previous tutorial, we
could search for:

.. code-block:: javascript

    rp({
        method: 'GET',
        uri: openrelayBaseURL + "/v0/orders?makerTokenAddress=0xe41d2489571d322189246dafa5ebde1f4699f498&takerTokenAddress=0x2956356cd2a2bf3202f771f50d3d14a367b48070",
        json: true,
    }).then((orders) => {
        for(var order of orders) {
            // Manipulate order object
        }
    })

Which would return orders with the same token pair as the order we issued in
the last tutorial.

Fill Orders
-----------

Once you have found an order you wish to fill as a taker, you need to submit
it to the 0x :term:`Exchange Contract`.

For this example, we will assume you have found an order offering 5000 UET for
0.1 WETH, with a taker fee of 0.25 ZRX. You don't want to fill the whole order,
you're only interested in buying 500 UET. You will pay 0.01 WETH and a fee of
0.025 ZRX to purchase 500 UET by filling this order. We will also assume that
you want to fill this order with the default account on your :term:`Web3`
object.

First, you have to have 0.01 WETH and 0.025 ZRX. Acquiring those tokens is
outside the scope of this tutorial.

Assuming that you have retrieved the order as a JSON object, you'll need to
convert some string values into BigNumber objects before the Exchange Contract
will accept them:

.. code-block:: javascript

    order.takerTokenAmount = new BigNumber(order.takerTokenAmount);
    order.makerTokenAmount = new BigNumber(order.makerTokenAmount);
    order.takerFee = new BigNumber(order.takerFee);
    order.makerFee = new BigNumber(order.makerFee);


Next you have to authorize the :term:`Token Proxy` to transfer those tokens on
your behalf.


.. code-block:: javascript

    var takerAllowance = zeroEx.token.setUnlimitedProxyAllowanceAsync(
        order.takerToken,
        web3.eth.accounts[0],
    );
    var feeAllowance = zeroEx.token.setUnlimitedProxyAllowanceAsync(
        "0xe41d2489571d322189246dafa5ebde1f4699f498",
        web3.eth.accounts[0],
    );

Once those transactions have completed, run:

.. code-block:: javascript

    var zeroEx = new ZeroEx.ZeroEx(web3.currentProvider);
    zeroEx.exchange.fillOrderAsync(
        signedOrder,
        // UET has 18 decimal places, so we need to provide the base unit
        new BigNumber("500000000000000000000"),
        false,
        web3.eth.accounts[0]
    )

If both you and the order maker have the necessary funds, and the order has not
already been filled, the funds will have been tranfered upon completion of this
transaction.
