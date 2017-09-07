import contextlib
import time

def bytesToInt(data):
    total = 0
    for i, byte in enumerate(data[::-1]):
        total += (256**i) * byte
    return total

# Specifically targeting uint256
def intToBytes(value):
    byteArray = []
    for i in range(32):
        position = value // (256**(31-i))
        byteArray.append(position)
        value -= position
    return bytes(byteArray)

def hexStringToBytes(value):
    if value.startswith("0x"):
        value = value[2:]
    result = []
    for i in range(len(value)//2):
        result.append(int(value[i*2:(i+1)*2], 16))
    return bytes(result)

def bytesToHexString(value):
    return "0x" + "".join(hex(b)[2:].zfill(2) for b in value)

class Locker(object):
    def __init__(self, redisClient):
        self.redisClient = redisClient

    @contextlib.contextmanager
    def lock(self, key, timeout=5):
        while True:
            if self.redisClient.setnx(key, str(int(time.time() + timeout))):
                try:
                    yield
                finally:
                    self.redisClient.delete(key)
                break
            else:
                if int(self.redisClient.get(key)) < time.time():
                    self.redisClient.delete(key)
                else:
                    time.sleep(0.250)

@contextlib.contextmanager
def get_queue_message(queue, redisClient):
    message = redisClient.brpoplpush(queue, "%s::unacked" % queue)
    try:
        yield message
    except Exception:
        redisClient.lpush("%s::error" % queue, message)
    redisClient.lrem("%s::unacked" % queue, 1, message)
