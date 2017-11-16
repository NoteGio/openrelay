#!/bin/sh

cd $(dirname $0)

REDIS_URL=${1:-redis://redis:6379}

ETHEREUM_URL="$ETHEREUM_URL" ./node_modules/.bin/truffle migrate --network main
ETHEREUM_URL="$ETHEREUM_URL" ./node_modules/.bin/truffle exec netinit.js $REDIS_URL --network main
