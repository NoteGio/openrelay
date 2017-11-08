import argparse
import logging
import json

import dynamo
import order
import util

logger = logging.getLogger(__name__)

def record_order(data, locker):
    order_obj = order.Order.FromBytes(data)
    # Make sure that only one process at a time is updating a given order
    with locker.lock(order_obj.orderHash):
        try:
            dynamo_order = dynamo.DynamoOrder.get(order_obj.orderHash)
        except dynamo.DynamoOrder.DoesNotExist:
            dynamo_order = dynamo.DynamoOrder.FromOrder(order_obj)
        # If the incoming record shows a higher takerTokenFilledAmount than
        # we've previously stored, update it
        stored_filled = util.bytesToInt(dynamo_order.takerTokenAmountFilled)
        if order_obj.takerTokenAmountFilled > stored_filled:
            incoming_filled = util.intToBytes(order_obj.takerTokenAmountFilled)
            dynamo_order.takerTokenAmountFilled = incoming_filled
        stored_cancelled = util.bytesToInt(dynamo_order.takerTokenAmountCancelled)
        if order_obj.takerTokenAmountCancelled > stored_cancelled:
            incoming_cancelled = util.intToBytes(order_obj.takerTokenAmountCancelled)
            dynamo_order.takerTokenAmountFilled = incoming_cancelled
        dynamo_order.save()


# def record_fill(orderHash, filled_amount, locker):
#     return dynamo.DynamoOrder.addFilled(orderHash, filled_amount, locker)


def delete_order(data, locker):
    order_obj = order.Order.FromBytes(data)
    with locker.lock(order_obj.orderHash):
        try:
            dynamo_order = dynamo.DynamoOrder.get(order_obj.orderHash)
        except dynamo.DynamoOrder.DoesNotExist:
            pass
        else:
            logger.info("Deleting order %s" % util.bytesToHexString(order_obj.orderHash))
            dynamo_order.delete()


def index_orders(redis_url, order_queue, unindex=False):
    redisClient = util.get_redis_client(redis_url)
    while True:
        with util.get_queue_message(order_queue, redisClient) as message:
            if not unindex:
                try:
                    record_order(message, util.Locker(redisClient))
                except Exception:
                    logger.exception("Error recording message")
            else:
                try:
                    delete_order(message, util.Locker(redisClient))
                except Exception:
                    logger.exception("Error deleting record")

# def fill_monitor():
#     parser = argparse.ArgumentParser()
#     parser.add_argument("redis_url")
#     parser.add_argument("fill_queue")
#     args = parser.parse_args
#
#     redisClient = get_redis_client(args.redis_url)
#
#     while True:
#         with util.get_queue_message(args.fill_queue, redisClient) as message:
#             fill = json.loads(message.decode("utf8"))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("redis_url")
    parser.add_argument("order_queue")
    parser.add_argument("--create", action="store_true", default=False)
    parser.add_argument("--unindex", action="store_true", default=False)
    args = parser.parse_args()
    if args.create and not dynamo.DynamoOrder.exists():
        dynamo.DynamoOrder.create_table(wait=True)
    index_orders(args.redis_url, args.order_queue, args.unindex)
