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
        projection = IncludeProjection(["data", "makerTokenAmountFilled"])
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
        projection = IncludeProjection(["data", "makerTokenAmountFilled"])
    makerToken = BinaryAttribute(hash_key=True)


class DynamoOrderTakerTokenIndex(GlobalSecondaryIndex):
    """
    Index to allow searching by takertoken
    """
    class Meta:
        index_name = 'order-takertoken-idx'
        read_capacity_units = 1
        write_capacity_units = 1
        projection = IncludeProjection(["data", "makerTokenAmountFilled"])
    takerToken = BinaryAttribute(hash_key=True)


class DynamoOrder(Model):
    """
    PynamoDB representation of an Order
    """
    class Meta:
        table_name = "Order"
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
    makerTokenAmountFilled = BinaryAttribute()
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
        self.makerTokenAmountFilled = util.intToBytes(order.makerTokenAmountFilled)
        self.pairHash = order.pairHash
        self.price = order.price
        self.data = order.rawdata
        return self

    def ToOrder(self):
        return order.Order.FromBytes(self.data)

    @classmethod
    def addFilled(cls, orderHash, amountFilled, locker):
        with locker.lock("%s::lock" % orderHash):
            order_dynamo = cls.get(orderHash)
            totalFilled = util.intToBytes(order_dynamo.makerTokenAmountFilled)
            totalFilled += amountFilled
            order_dynamo.makerTokenAmountFilled = util.bytesToInt(totalFilled)
            order_dynamo.save()
