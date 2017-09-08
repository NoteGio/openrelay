import unittest
import util
import order
import sample

class OrderTestCase(unittest.TestCase):
    def test_order_from_bytes(self):
        test_order = order.Order.FromBytes(sample.data)
        self.assertEqual(test_order.price, 0.02)
        self.assertEqual(test_order.orderHash, util.hexStringToBytes(
            "731319211689ccf0327911a0126b0af0854570c1b6cdfeb837b0127e29fe9fd5"
        ))
    def test_to_dict(self):
        test_order = order.Order.FromBytes(sample.data)
        order_dict = test_order.to_dict()

        self.assertEqual(order_dict["makerTokenAddress"], "0x1dad4783cf3fe3085c1426157ab175a6119a04ba")
        self.assertEqual(order_dict["maker"], "0x324454186bb728a3ea55750e0618ff1b18ce6cf8")
        self.assertEqual(order_dict["taker"], "0x0000000000000000000000000000000000000000")
        self.assertEqual(order_dict["feeRecipient"], "0x0000000000000000000000000000000000000000")
        self.assertEqual(order_dict["takerTokenAddress"], "0x05d090b51c40b020eab3bfcb6a2dff130df22e9c")
        self.assertEqual(order_dict["exchangeContractAddress"], "0x90fe2af704b34e0224bf2299c838e04d4dcf1364")
        self.assertEqual(order_dict["makerTokenAmount"], "50000000000000000000")
        self.assertEqual(order_dict["takerTokenAmount"], "1000000000000000000")
        self.assertEqual(order_dict["makerFee"], "0")
        self.assertEqual(order_dict["takerFee"], "0")
        self.assertEqual(order_dict["expirationUnixTimestampSec"], "1502841540")
        self.assertEqual(order_dict["salt"], "11065671350908846865864045738088581419204014210814002044381812654087807531")
        self.assertEqual(order_dict["ecSignature"]["v"], 27)
        self.assertEqual(order_dict["ecSignature"]["r"], "0x021fe6dba378a347ea5c581adcd0e0e454e9245703d197075f5d037d0935ac2e")
        self.assertEqual(order_dict["ecSignature"]["s"], "0x12ac107cb04be663f542394832bbcb348deda8b5aa393a97a4cc3139501007f1")
