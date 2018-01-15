#!/bin/sh

REDIS_HOST=$1
PREWARM=$2

geth --syncmode=light &

GETH_PID=$!

echo $GETH_PID

if [ -z "$REDIS_HOST" ]; then
  BLOCKNUMBER=$(echo "get topic://newblocks::blocknumber" | nc $REDIS_HOST 6379 | tr -d '\r')
fi

sleep 5

while kill -0 $GETH_PID
do
    while geth --exec "if(admin.peers.length >= 1 && eth.blockNumber > ${BLOCKNUMBER:-4883209} && "'!'"eth.syncing){admin.startRPC('0.0.0.0'); } else { console.log('notready')}" attach | grep notready
    do
        sleep 5
    done
    if ! [ -z "$PREWARM" ]; then
        echo "Prewarm flag set. Terminating"
        exit 0
    fi
    PEERCOUNT=$(geth --exec "console.log(admin.peers.length)" attach | grep -v undefined)
    while [ "$PEERCOUNT" -ge "1" ]
    do
        PEERCOUNT=$(geth --exec "console.log(admin.peers.length)" attach | grep -v undefined)
        BLOCKNUMBER=$(geth --exec "console.log(eth.blockNumber)" attach | grep -v undefined)
        sleep 10
    done
    geth --exec "admin.stopRPC()" attach
    echo "Lost all peers. Stopped serving RPC. Waiting for new peers."
done
