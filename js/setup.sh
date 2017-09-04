#!/bin/sh

cd $(dirname $0)

REDIS_URL=${1:-redis://redis:6379}

./node_modules/.bin/truffle migrate --network testnet
./node_modules/.bin/truffle exec netinit.js $REDIS_URL --network testnet
