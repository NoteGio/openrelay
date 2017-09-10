import dynamo
import util
import itertools
import json
import hashlib

from flask import Flask, request, make_response

app = Flask(__name__)

def format_response(orders, count, accept):
    items = itertools.islice(orders, count)
    if "application/octet-stream" in accept:
        resp = make_response(b"".join(item.binary() for item in items))
        resp.headers["Content-Type"] = "application/octet-stream"
    else:
        resp = make_response(json.dumps({
            "orders": [item.ToOrder().to_dict() for item in items]
        }))
        resp.headers["Content-Type"] = "application/json"
    return resp

@app.route('/mtok/<maker_token>')
def maker_token_search(maker_token):
    orders = dynamo.DynamoOrder.makertoken_index.query(
        util.hexStringToBytes(maker_token)
    )
    return format_response(
        orders,
        request.args.get("count", 25),
        request.headers.get("Accept", "")
    )

@app.route('/ttok/<taker_token>')
def taker_token_search(taker_token):
    orders = dynamo.DynamoOrder.takertoken_index.query(
        util.hexStringToBytes(taker_token)
    )
    return format_response(
        orders,
        request.args.get("count", 25),
        request.headers.get("Accept", "")
    )

@app.route('/<maker_token>/<taker_token>')
def pair_search(maker_token, taker_token):

    orders = dynamo.DynamoOrder.pairhash_index.query(
        hashlib.sha256(
            util.hexStringToBytes(maker_token) +
            util.hexStringToBytes(taker_token)
        ).digest()
    )
    return format_response(
        orders,
        request.args.get("count", 25),
        request.headers.get("Accept", "")
    )
