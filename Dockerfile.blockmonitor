FROM python:3.6

ADD js/blockmonitor.py /project/blockmonitor.py

RUN pip install requests redis

WORKDIR /project

CMD ["python", "blockmonitor.py", "http://ethnode:8545", "redis:6379", "queue://newblocks"]
