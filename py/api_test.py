import json
import unittest
import itertools

import sample
import order
import dynamo
import flask
import api
import util

class ApiTestCase(unittest.TestCase):
    def setUp(self):
        self.orders = [
            dynamo.DynamoOrder.FromOrder(order.Order.FromBytes(s))
            for s in sample.samples
        ]
        self.app = flask.Flask("test")

    def test_format_orders_json(self):
        with self.app.app_context():
            res = api.format_response(self.orders, 25, "application/json")
            data = json.loads(res.data.decode("utf8"))
            for order_dict, dynamo_order in zip(data["orders"], self.orders):
                self.assertEqual(
                    order_dict["makerTokenAddress"],
                    util.bytesToHexString(dynamo_order.makerToken)
                )
                self.assertEqual(
                    order_dict["takerTokenAddress"],
                    util.bytesToHexString(dynamo_order.takerToken)
                )
                self.assertEqual(
                    util.intToBytes(int(order_dict["makerTokenAmountFilled"])),
                    dynamo_order.makerTokenAmountFilled
                )

    def test_format_orders_bin(self):
        with self.app.app_context():
            res = api.format_response(self.orders, 25, "application/octet-stream")
            bin_orders = [res.data[409*i:409*(i+1)]
                          for i in range(len(res.data) // 409)]

            for bin_order, dynamo_order in zip(bin_orders, self.orders):
                cmp_order = order.Order.FromBytes(bin_order)
                self.assertEqual(
                    cmp_order.orderHash, dynamo_order.orderHash
                )
