import dynamo
import util
import itertools
import json

from flask import Flask, request, make_response

app = Flask(__name__)

@app.route('/maker/<maker_token>')
def maker_search(maker_token):
    orders = dynamo.DynamoOrder.makertoken_index.query(
        util.hexStringToBytes(maker_token)
    )
    items = itertools.islice(orders, request.args.get("count", 25))
    if "application/octet-stream" in request.headers.get("Accept", ""):
        resp = make_response(b"".join(item.binary() for item in items))
        resp.headers["Content-Type"] = "application/octet-stream"
    else:
        resp = make_response(json.dumps([item.ToOrder().to_dict() for item in items]))
        resp.headers["Content-Type"] = "application/json"
    return resp
