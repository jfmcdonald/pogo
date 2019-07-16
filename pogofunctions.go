package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

// this code taken from https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
// setup data needed to generate random stirngs
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	// set the rand seed
	var src = rand.NewSource(time.Now().UnixNano())
	// create a sice to hold our byte string
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// End borrowed code.

// Description: get a file path from redis
// Expects: pointer to redis Client
// Retruns a file path as String
// TODO:
func getPath(client *redis.Client) (string, string) {
	// get a random file key from the db:s
	key, randErr := client.RandomKey().Result()
	if randErr != nil {
		log.Println("unable to get randome key:")
		log.Println(randErr)
	}
	// get the keys value
	val, getErr := client.Get(key).Result()
	if getErr != nil {
		log.Println("unable to get value from key:")
		log.Println(getErr)
	}
	return key, val
}

// Description: get a count of the number of keys in a redis set
// Expects: pointer to open redis client
// Returns: uint64 number of keys
func getKeysCount(client *redis.Client) uint64 {
	keys, countErr := client.DBSize().Result()
	if countErr != nil {
		log.Println("getKeysCount: unalbe to get count: ")
		log.Println(countErr)
		panic(countErr)
	}
	// returning a bit memory space since there could be a lot of files
	return uint64(keys)
}
