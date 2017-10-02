import time
import argparse
import logging

import dynamo

import util

logger = logging.getLogger(__name__)


def queue_orders(redisClient, publish_queue, queued_max, length_check_frequency=0.1):
    queued = 0
    counter = 0
    for order in dynamo.DynamoOrder.scan():
        if order.ToOrder().expirationTimestampInSec < time.time():
            # We can prune expired orders as we go without queueing
            order.delete()
            continue
        redisClient.lpush(publish_queue, order.data)
        queued += 1
        counter += 1
        while queued >= queued_max:
            # If we have millions of orders in dynamo, we don't want to swamp
            # redis with all of them at once. Once we hit the maximum number of
            # items in the queue, don't add more until some have been consumed

            # This should take some pressure off of both Redis and Dynamo.
            queued = redisClient.llen(publish_queue)
            if queued >= queued_max:
                time.sleep(length_check_frequency)

    logger.info("Queued %s items" % counter)


def main(redisClient, publish_queue, delay, queued_max, length_check_frequency=0.1):
    while True:
        start_time = time.time()
        queue_orders(redisClient, publish_queue, queued_max,
                     length_check_frequency)
        while redisClient.llen(publish_queue) > 0:
            # If we're doing this repeatedly, we wait until the queue is
            # cleared to load it up again.
            time.sleep(length_check_frequency)
        remaining_time = (start_time + delay) - time.time()
        if remaining_time > 0:
            # If the queue clears faster than the round delay, wait until the
            # time is up to load the queue again.
            time.sleep(remaining_time)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("redis_url")
    parser.add_argument("publish_queue")
    parser.add_argument("--round-delay", "-d", type=int, default=60)
    parser.add_argument("--length-check-freq", "-f", type=int, default=0.1)
    parser.add_argument("--recur", "-r", action="store_true", default=False)
    parser.add_argument("--queued-max", "-m", type=int, default=1000)
    parser.add_argument("--log-level", "-l", default="info")
    args = parser.parse_args()

    logging.basicConfig(level=getattr(logging, args.log_level.upper()))

    redisClient = util.get_redis_client(args.redis_url)
    if args.recur:
        main(
            redisClient,
            args.publish_queue,
            args.round_delay,
            args.queued_max,
            args.length_check_freq
        )
    else:
        queue_orders(
            redisClient,
            args.publish_queue,
            args.queued_max,
            args.length_check_freq
        )
