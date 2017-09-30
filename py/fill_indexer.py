import util
import argparse
import json
import dynamo

import logging

logger = logging.getLogger(__name__)


def process_fill(fill, locker):
    orderHash = util.hexStringToBytes(fill["orderHash"])
    dynamo_order = dynamo.DynamoOrder.addFilled(
        orderHash,
        int(fill.get("filledTakerTokenAmount", 0)),
        int(fill.get("cancelledTakerTokenAmount", 0)),
        locker
    )
    order = dynamo_order.ToOrder()
    total_unavailable = (
        util.bytesToInt(dynamo_order.takerTokenAmountFilled) +
        util.bytesToInt(dynamo_order.takerTokenAmountCancelled)
    )
    if total_unavailable >= (order.takerTokenAmount * .99):
        # The order is > 99% filled, delist it.
        dynamo_order.delete()


def fill_monitor(redisClient, fill_queue, locker):
    while True:
        with util.get_queue_message(fill_queue, redisClient) as message:
            fill = json.loads(message.decode("utf8"))
            logger.debug("Updated: %s" % message)
            process_fill(fill, locker)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("redis_url")
    parser.add_argument("fill_queue")
    parser.add_argument("--log-level", "-l", default="info")
    args = parser.parse_args()

    logging.basicConfig(level=getattr(logging, args.log_level.upper()))

    redisClient = util.get_redis_client(args.redis_url)
    locker = util.Locker(redisClient)
    fill_monitor(redisClient, args.fill_queue, locker)
