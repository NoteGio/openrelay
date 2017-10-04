import dynamo
import util
import itertools
import json
import hashlib
import threading
import argparse
from functools import wraps

from flask import Flask, request, make_response, redirect
from flask_cors import CORS

app = Flask(__name__)
CORS(app)
blockhash = "pending"


def format_response(orders, count, accept):
    items = itertools.islice(orders, int(count))
    if "application/octet-stream" in accept:
        resp = make_response(b"".join(item.binary() for item in items))
        resp.headers["Content-Type"] = "application/octet-stream"
    else:
        resp = make_response(json.dumps([
            item.ToOrder().to_dict() for item in items
        ]))
        resp.headers["Content-Type"] = "application/json"
    return resp

def notsupported(field):
    return make_response(json.dumps({
        "err": ("The search field '%s' is not currently supported by OpenRelay"
                % field)
    }), 501)

UNSUPPORTED_FIELDS = [
    "ascByBaseToken",
    "exchangeContractAddress",
    "isExpired",
    "isOpen",
    "isClosed",
    "token",
    "maker",
    "taker",
    "trader",
    "feeRecipient",
]

def req_blockhash(fn):
    @wraps(fn)
    def wrapper(*args, **kwargs):
        if request.args.get("blockhash", None) is None:
            if "?" in request.full_path.rstrip("?"):
                new_path = "%s&blockhash=%s" % (
                    request.full_path.rstrip("?"),
                    blockhash.decode("utf8").strip('"')
                )
            else:
                new_path = "%s?blockhash=%s" % (
                    request.full_path.rstrip("?"),
                    blockhash.decode("utf8").strip('"')
                )
            return redirect(new_path, 307)
        else:
            return fn(*args, **kwargs)
    return wrapper

@app.route('/')
def root():
    return redirect('/v0/orders')

@app.route('/v0/orders')
@req_blockhash
def orders():
    for field in UNSUPPORTED_FIELDS:
        if request.args.get(field, False):
            return notsupported(field)
    count = int(request.args.get("limit", 25))
    makerTokenAddress = request.args.get("makerTokenAddress", None)
    takerTokenAddress = request.args.get("takerTokenAddress", None)
    if makerTokenAddress and takerTokenAddress:
        lookp_hash = hashlib.sha256(
            util.hexStringToBytes(makerTokenAddress) +
            util.hexStringToBytes(takerTokenAddress)
        ).digest()
        index = dynamo.DynamoOrder.pairhash_index
    elif makerTokenAddress and not takerTokenAddress:
        lookup_hash = util.hexStringToBytes(makerTokenAddress)
        index = dynamo.DynamoOrder.makertoken_index
    elif not makerTokenAddress and takerTokenAddress:
        lookup_hash = util.hexStringToBytes(takerTokenAddress)
        index = dynamo.DynamoOrder.takertoken_index
    else:
        orders = dynamo.DynamoOrder.scan(limit=count)
        index = None
    if index:
        orders = index.query(lookup_hash, limit=count)

    return format_response(
        orders,
        count,
        request.headers.get("Accept", "")
    )

@app.route('/v0/order/<order_hash>')
@req_blockhash
def single_order(order_hash):
    order = dynamo.DynamoOrder.get(util.hexStringToBytes(order_hash))
    accept = request.headers.get("Accept", "")
    if "application/octet-stream" in accept:
        resp = make_response(order.binary())
        resp.headers["Content-Type"] = "application/octet-stream"
    else:
        resp = make_response(json.dumps(order.ToOrder().to_dict()))
        resp.headers["Content-Type"] = "application/json"
    return resp


@app.route('/_hc')
def health_check():
    orders = dynamo.DynamoOrder.scan(limit=1)
    return format_response(orders, 0, "")

@app.route('/mtok/<maker_token>')
@req_blockhash
def maker_token_search(maker_token):
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.makertoken_index.query(
        util.hexStringToBytes(maker_token),
        limit=count
    )
    return format_response(orders, count, request.headers.get("Accept", ""))

@app.route('/ttok/<taker_token>')
@req_blockhash
def taker_token_search(taker_token):
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.takertoken_index.query(
        util.hexStringToBytes(taker_token),
        limit=count
    )
    return format_response(orders, count, request.headers.get("Accept", ""))

@app.route('/<maker_token>/<taker_token>')
@req_blockhash
def pair_search(maker_token, taker_token):
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.pairhash_index.query(
        hashlib.sha256(
            util.hexStringToBytes(maker_token) +
            util.hexStringToBytes(taker_token)
        ).digest(),
        scan_index_forward=(request.args.get("asc", "true") == "true"),
        limit=count
    )
    return format_response(orders, count, request.headers.get("Accept", ""))

def populate_blockhash(redis_url, topic_name):
    global blockhash
    r = util.get_redis_client(redis_url)
    # Get block number initially. The pubsub channel will give us block hashes,
    # but from a caching perspective they should all be unique, which is the
    # important part.
    blockhash = r.get("topic://%s::blocknumber" % topic_name)
    p = r.pubsub()
    p.subscribe(topic_name)
    for message in p.listen():
        if message.get("type") == "message":
            try:
                blockhash = message["data"]
            except KeyError:
                pass

if __name__ == "__main__":
    import bjoern
    parser = argparse.ArgumentParser()
    parser.add_argument("redis_url")
    parser.add_argument("topic_name")
    parser.add_argument("--port", type=int, default=8888)
    args = parser.parse_args()

    blockhashThread = threading.Thread(
        target=populate_blockhash, args=(args.redis_url, args.topic_name))
    blockhashThread.start()
    bjoern.run(app, "0.0.0.0", args.port)
