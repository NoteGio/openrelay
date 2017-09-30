import os
import unittest
import contextlib

import dynamo
import indexer
import fill_indexer
import sample
import util


class Locker(object):
    @contextlib.contextmanager
    def lock(self, name):
       yield


class IndexerTestCase(unittest.TestCase):
    def setUp(self):
        try:
            os.environ["DYNAMODB_HOST"]
        except KeyError:
            raise self.skipTest("No dynamodb configuration")
        dynamo.DynamoOrder.create_table(wait=True)

    def tearDown(self):
        dynamo.DynamoOrder.delete_table()

    def test_index(self):
        indexer.record_order(sample.data, Locker())
        item = next(dynamo.DynamoOrder.scan())
        test_order = item.ToOrder()
        self.assertEqual(test_order.price, 0.02)
        self.assertEqual(test_order.orderHash, util.hexStringToBytes(
            "731319211689ccf0327911a0126b0af0854570c1b6cdfeb837b0127e29fe9fd5"
        ))

    def test_fill(self):
        indexer.record_order(sample.data, Locker())
        h_ = "731319211689ccf0327911a0126b0af0854570c1b6cdfeb837b0127e29fe9fd5"
        fill_indexer.process_fill({
            "orderHash": h_,
            "filledMakerTokenAmount": 25000000000000000000
        }, Locker())
        item = next(dynamo.DynamoOrder.scan())
        self.assertEqual(
            util.bytesToInt(item.takerTokenAmountFilled),
            25000000000000000000
        )
        fill_indexer.process_fill({
            "orderHash": h_,
            "filledMakerTokenAmount": 25000000000000000000
        }, Locker())
        with self.assertRaises(StopIteration):
            next(dynamo.DynamoOrder.scan())
