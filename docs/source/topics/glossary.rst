Glossary
========

.. glossary::

    Order
    Offer
        An order, sometimes called an offer, is a request to trade a certain
        amount of one token for a certain amount of another token. An order is
        cryptographically signed to prove that the account that made the order
        intends for it to be executed.

    Maker
        The account that creates and signs an :term:`order`. Within an order,
        the maker is indicated by an Ethereum account address.

    Taker
        The account that fills an :term:`order` with the 0x Exchange Contract.
        Sometimes, an order may specify a taker address, in which case that is
        the only account that may :term:`fill` the order. Other times the taker
        address is left blank, allowing anyone to fill the order.

    Maker Token
        The token that the :term:`maker` of an :term:`order` is offering.

    Taker Token
        The token that the :term:`taker` of an :term:`order` must provide to
        the :term:`maker` to complete the order.

    Ethereum
        `Ethereum <https://ethereum.org/>`_ is a blockchain that allows the
        execution of general purpose smart contracts. Each token that can be
        traded on OpenRelay is its own Ethereum smart contract. The 0x
        :term:`Exchange` is also an Ethereum smart contract.

    0x Protocol
        The 0x protocol allows users to trustlessly exchange
        :ref:`ERC20 tokens`. It uses off-chain :ref:`relays`, and settles
        transactions through the 0x :ref:`Exchange Contract`.

    Exchange
    Exchange Contract
        The 0x Exchange Contract allows Ethereum accounts to :term:`fill`
        orders. It verifies the cryptographic signature on the :ref:`order`,
        and transfers tokens between the :ref:`maker` and the :ref:`taker`, as
        well as transferring any fees to the :ref:`relay`.

    ERC20 Token
        The `ERC20 Token Standard <https://theethereum.wiki/w/index.php/ERC20_Token_Standard>`_
        defines a specific set of methods and events that an Ethereum contract
        should implement to have a transferrable token. To be tradeable via the
        :ref:`0x Protocol`, tokens must implement the ERC20 specification.

    WETH
        The Ethereum base currency, Ether, exists independently of any single
        contract, and does not comply with the ERC20 specification. To be able
        to trade Ether through 0x, which only works with ERC20 tokens, we must
        convert it into an :term:`ERC20 token`. WETH is an ERC20 token
        contract. You can send Ether to the WETH contract, and it will issue
        you ERC20 tokens (WETH) equal to the amount of Ether sent. At any time
        you may redeem your WETH tokens from the contract, and it will replace
        them with an equal number of Ether.

    Order Book
        An order book is a collection of :term:`orders`. OpenRelay provides an
        API for searching through all open orders.

    Fill
        An :term:`order` is filled when a :term:`taker` submits the order to
        the 0x :term:`Exchange Contract`. The taker is responsible for calling
        one of the contract's fill functions, and is must pay the :ref:`gas`
        for the Ethereum transaction.

    Gas
        `Gas <https://ethereum.stackexchange.com/a/62/18148>`_ is a proxy for
        the computational requirements of an Ethereum transaction. Every
        Ethereum transaction has a gas limit (a maximum amount of computation
        allowed by the transaction). When a transaction is submitted, it
        specifies a gas price, indicating the amount of Ether the submitter is
        willing to pay for each unit of gas. In general, transactions with high
        gas prices get included into blocks sooner than transactions with low
        gas prices.

    Token Proxy
        The 0x Protocol uses a contract called a Token Proxy to transfer tokens
        between the :ref:`maker` and :ref:`taker` when an :ref:`order` is
        :ref:`filled`. Users must approve the Token Proxy to transfer
        :ref:`ERC20 Tokens` on behalf of the user.

    Web3
        Web3 provides a standard interface for interacting with
        :term:`Ethereum` in JavaScript. There are several different Web3
        providers, including browsers like `Metamask <https://metamask.io/>`_,
        and `Mist <https://github.com/ethereum/mist>`_, as well as a
        `Node JS Package <https://www.npmjs.com/package/web3>`_.
