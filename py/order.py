import hashlib
import sha3


import util

class Order(object):
    @classmethod
    def FromBytes(cls, data):
        self = cls()
        self.rawdata = data
        self.orderHash = self.digest()
        self.exchangeAddress = data[0:20]
        self.maker = data[20:40]
        self.taker = data[40:60]
        self.makerToken = data[60:80]
        self.takerToken = data[80:100]
        self.feeRecipient = data[100:120]
        self.makerTokenAmount = util.bytesToInt(data[120:152])
        self.takerTokenAmount = util.bytesToInt(data[152:184])
        self.makerFee = util.bytesToInt(data[184:216])
        self.takerFee = util.bytesToInt(data[216:248])
        self.expirationTimestampInSec = util.bytesToInt(data[248:280])
        self.salt = util.bytesToInt(data[280:312])
        self.sigV = data[312]
        self.sigR = data[313:345]
        self.sigS = data[345:377]
        self.makerTokenAmountFilled = util.bytesToInt(data[377:409])
        self.price = self.takerTokenAmount / self.makerTokenAmount
        self.pairHash = hashlib.sha256(self.makerToken + self.takerToken).digest()
        return self

    def digest(self):
        orderHash = sha3.keccak_256()
        orderHash.update(self.rawdata[:312])
        return orderHash.digest()

    def to_dict(self):
        return {
          "expirationUnixTimestampSec": str(self.expirationTimestampInSec),
          "feeRecipient": util.bytesToHexString(self.feeRecipient),
          "maker": util.bytesToHexString(self.maker),
          "makerFee": str(self.makerFee),
          "makerTokenAmount": str(self.makerTokenAmount),
          "makerTokenAmountFilled": str(self.makerTokenAmountFilled),
          "salt": str(self.salt),
          "taker": util.bytesToHexString(self.taker),
          "takerFee": str(self.takerFee),
          "takerTokenAmount": str(self.takerTokenAmount),
          "makerTokenAddress": util.bytesToHexString(self.makerToken),
          "takerTokenAddress": util.bytesToHexString(self.takerToken),
          "exchangeContractAddress": util.bytesToHexString(self.exchangeAddress),
          "ecSignature": {
            "v": self.sigV,
            "r": util.bytesToHexString(self.sigR),
            "s": util.bytesToHexString(self.sigS)
          }
        }
