package channels

import (
	"gopkg.in/redis.v3"
	"log"
	"strings"
)

func redisErrIsNil(result redis.Cmder) bool {
	switch result.Err() {
	case nil:
		return false
	case redis.Nil:
		return true
	default:
		if strings.HasSuffix(result.Err().Error(), "i/o timeout") {
			return true
		}
		log.Printf("rmq redis error is not nil %s", result.Err())
		return false
	}
}
