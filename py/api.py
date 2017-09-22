import dynamo
import util
import itertools
import json
import hashlib

from flask import Flask, request, make_response

app = Flask(__name__)

# TODO: require blocknumber argument

def format_response(orders, count, accept):
    items = itertools.islice(orders, int(count))
    if "application/octet-stream" in accept:
        resp = make_response(b"".join(item.binary() for item in items))
        resp.headers["Content-Type"] = "application/octet-stream"
    else:
        resp = make_response(json.dumps({
            "orders": [item.ToOrder().to_dict() for item in items]
        }))
        resp.headers["Content-Type"] = "application/json"
    return resp

@app.route('/')
def scan_all():
    count = int(request.args.get("count", 25))
    orders = dynamo.DynamoOrder.scan(limit=count)
    return format_response(
        orders,
        count,
        request.headers.get("Accept", "")
    )

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
