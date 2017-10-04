import dynamo
import util
import itertools
import json
import hashlib

from flask import Flask, request, make_response, redirect
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

# TODO: require blocknumber argument

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

@app.route('/')
def root():
    return redirect('/v0/orders')

@app.route('/v0/orders')
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
def maker_token_search(maker_token):
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.makertoken_index.query(
        util.hexStringToBytes(maker_token),
        limit=count
    )
    return format_response(orders, count, request.headers.get("Accept", ""))

@app.route('/ttok/<taker_token>')
def taker_token_search(taker_token):
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.takertoken_index.query(
        util.hexStringToBytes(taker_token),
        limit=count
    )
    return format_response(orders, count, request.headers.get("Accept", ""))

@app.route('/<maker_token>/<taker_token>')
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

if __name__ == "__main__":
    import bjoern
    bjoern.run(app, "0.0.0.0", 8888)
