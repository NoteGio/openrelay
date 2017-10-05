#!/bin/sh

geth --syncmode=light --rpcaddr 0.0.0.0 &

sleep 5

while geth --exec 'if(eth.blockNumber > 4339872 && !eth.syncing){admin.startRPC("0.0.0.0"); } else { console.log("notready")}' attach | grep notready
do
    sleep 1
done

wait
