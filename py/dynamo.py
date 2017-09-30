import os

from pynamodb.models import Model
from pynamodb.indexes import GlobalSecondaryIndex, IncludeProjection
from pynamodb.attributes import (
    BinaryAttribute, NumberAttribute, UTCDateTimeAttribute
)
import order
import util


class DynamoOrderPairhashIndex(GlobalSecondaryIndex):
    """
    Index to allow searching by makertoken / takertoken pairs, ordering by
    price
    """
    class Meta:
        index_name = 'order-pairhash-idx'
        read_capacity_units = 1
        write_capacity_units = 1
        projection = IncludeProjection([
                                        "data",
                                        "takerTokenAmountFilled",
                                        "takerTokenAmountCancelled"
                                        ])
    pairHash = BinaryAttribute(hash_key=True)
    price = NumberAttribute(range_key=True)


class DynamoOrderMakerTokenIndex(GlobalSecondaryIndex):
    """
    Index to allow searching by makertoken
    """
    class Meta:
        index_name = 'order-makertoken-idx'
        read_capacity_units = 1
        write_capacity_units = 1
        projection = IncludeProjection([
                                        "data",
                                        "takerTokenAmountFilled",
                                        "takerTokenAmountCancelled"
                                        ])
    makerToken = BinaryAttribute(hash_key=True)


class DynamoOrderTakerTokenIndex(GlobalSecondaryIndex):
    """
    Index to allow searching by takertoken
    """
    class Meta:
        index_name = 'order-takertoken-idx'
        read_capacity_units = 1
        write_capacity_units = 1
        projection = IncludeProjection([
                                        "data",
                                        "takerTokenAmountFilled",
                                        "takerTokenAmountCancelled"
                                        ])
    takerToken = BinaryAttribute(hash_key=True)


class DynamoOrder(Model):
    """
    PynamoDB representation of an Order
    """
    class Meta:
        table_name = os.environ.get("ORDER_TABLE_NAME", "Order")
        read_capacity_units = 1
        write_capacity_units = 1
        region = os.environ.get("AWS_REGION", "us-east-2")
        try:
            host = os.environ["DYNAMODB_HOST"]
        except KeyError:
            pass
    orderHash = BinaryAttribute(hash_key=True)
    makerToken = BinaryAttribute()
    takerToken = BinaryAttribute()
    takerTokenAmountFilled = BinaryAttribute()
    takerTokenAmountCancelled = BinaryAttribute()
    pairHash = BinaryAttribute()
    price = NumberAttribute()
    data = BinaryAttribute()
    pairhash_index = DynamoOrderPairhashIndex()
    makertoken_index = DynamoOrderMakerTokenIndex()
    takertoken_index = DynamoOrderTakerTokenIndex()

    @classmethod
    def FromOrder(cls, order):
        self = cls()
        self.orderHash = order.orderHash
        self.makerToken = order.makerToken
        self.takerToken = order.takerToken
        self.takerTokenAmountFilled = util.intToBytes(order.takerTokenAmountFilled)
        self.takerTokenAmountCancelled = util.intToBytes(order.takerTokenAmountCancelled)
        self.pairHash = order.pairHash
        self.price = order.price
        self.data = order.rawdata[:377]
        return self

    def ToOrder(self):
        return order.Order.FromBytes(self.binary())

    def binary(self):
        return self.data + self.takerTokenAmountFilled + self.takerTokenAmountCancelled

    @classmethod
    def addFilled(cls, orderHash, amountFilled, amountCancelled, locker):
        with locker.lock("%s::lock" % orderHash):
            order_dynamo = cls.get(orderHash)
            totalFilled = util.bytesToInt(order_dynamo.takerTokenAmountFilled)
            totalFilled += amountFilled
            order_dynamo.takerTokenAmountFilled = util.intToBytes(totalFilled)
            totalCancelled = util.bytesToInt(order_dynamo.takerTokenAmountCancelled)
            totalCancelled += amountCancelled
            order_dynamo.takerTokenAmountCancelled = util.intToBytes(totalCancelled)
            order_dynamo.save()
            return order_dynamo
