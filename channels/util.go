package channels

import (
	"gopkg.in/redis.v3"
	"log"
)

func redisErrIsNil(result redis.Cmder) bool {
	switch result.Err() {
	case nil:
		return false
	case redis.Nil:
		return true
	default:
		log.Printf("rmq redis error is not nil %s", result.Err())
		return false
	}
}
