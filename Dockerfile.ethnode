FROM ethereum/client-go

ADD docker-cfg/run-eth.sh /run-eth.sh

RUN chmod +x /run-eth.sh

RUN /run-eth.sh "" --prewarm

RUN rm -f /root/.ethereum/nodekey

ENTRYPOINT ["/run-eth.sh"]
