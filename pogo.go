package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
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

func RandString(n int) string {
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

// Description: creates a new blank file and logs it into the redis db
// Expects: a pointers to the root path, key ttl and open redis client
// Returns: nothing
// TODO: currnetly panics on creation or connection error. this should be imporved
func createFile(pathPtr *string, ttlPtr *int, client *redis.Client) {
	// generate a random file name
	fileName := RandString(10)
	// build up the full path
	filePath := *pathPtr + "/" + fileName
	// Create the file on disk
	emptyFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	emptyFile.Close()
	// add the place holder to the redis db
	setErr := client.Set(fileName, filePath, time.Duration(*ttlPtr)*time.Second).Err()
	if setErr != nil {
		panic(setErr)
	}
}

// Description: updates an random existing file with 8 bits of data
// Expects: pointer to an open redis client
// Returns: nothing
// TODO: replace panics with proper logging
func updateFile(client *redis.Client) {
	// get a random key from redis
	key, randErr := client.RandomKey().Result()
	if randErr != nil {
		panic(randErr)
	}
	// get the keys value
	val, getErr := client.Get(key).Result()
	if getErr != nil {
		panic(getErr)
	}
	// open our file for writing, no read we don't want to create an
	// excuse for the filesystem to cache data
	file, ioErr := os.OpenFile(val, os.O_WRONLY|os.O_APPEND, 0644)
	if ioErr != nil {
		panic(ioErr)
	}
	// defer the close till the end of the funciton
	defer file.Close()
	// write out an random 8 bits. We are igoring the length here since
	// we don't use it
	_, writeErr := file.WriteString(RandString(8))
	if writeErr != nil {
		panic(writeErr)
	}
}

// Description: read a file by line, generate a line count. This is just to make
// sure the file is actually read from disk in its entirety
// Expects: pointer to open redis db connection
// Returns: nothing
// TODO: proper error reporting
func readFile(client *redis.Client) {
	// get a random file key from the db
	key, randErr := client.RandomKey().Result()
	if randErr != nil {
		panic(randErr)
	}
	// get the keys value
	val, getErr := client.Get(key).Result()
	if getErr != nil {
		panic(getErr)
	}
	// Read in our file
	data, readErr := ioutil.ReadFile(val)
	if readErr != nil {
		panic(readErr)
	}
	// do something with the data so it's really in mem
	data = data
}

// Description: Delete a random file and remove it's db entery
// Expects: pointer to open redis client
// Returns: nothing
// TODO: proper error logging, redis key delete error check is broken
func delFile(client *redis.Client) {
	// get a random file key from the db
	key, randErr := client.RandomKey().Result()
	if randErr != nil {
		panic(randErr)
	}
	// get the keys value
	val, getErr := client.Get(key).Result()
	if getErr != nil {
		panic(getErr)
	}
	remErr := os.Remove(val)
	if remErr != nil {
		panic(remErr)
	}
	client.Del(key)
	// This always errors even when it works
	//if redisDelErr != nil {
	//	panic(redisDelErr)
	//}
}

// Description: get a count of the number of keys in a redis set
// Expects: pointer to open redis client
// Returns: uint64 number of keys
func getKeysCount(client *redis.Client) uint64 {
	keys, countErr := client.DBSize().Result()
	if countErr != nil {
		panic(countErr)
	}
	return uint64(keys)
}

// Description: delete all remaining keys from the current run
// Expects: pointer to open redis client
// Returns: nothing
func delAllFiles(client *redis.Client) {
	keysCount := getKeysCount(client)
	for keysCount > 0 {
		key, randErr := client.RandomKey().Result()
		if randErr != nil {
			panic(randErr)
		}
		val, keyErr := client.Get(key).Result()
		if keyErr != nil {
			panic(keyErr)
		}
		remErr := os.Remove(val)
		if remErr != nil {
			panic(remErr)
		}
		client.Del(key)
		keysCount = getKeysCount(client)
	}
}
func main() {

	// command line options
	pathPtr := flag.String("path", "/tmp/pogo", "Path wher run time files will be generated")
	filecountPtr := flag.Uint64("count", 10, "Total number of files to generate")
	ttlPtr := flag.Int("ttl", 60, "Index Key/Value store default key TTL"
	redishostPtr := flag.String("dbhost", "localhost", "Hostname of the network redis server")
	redisdbPtr := flag.Int("db", 0, "redis db id you want to store keys in")

	flag.Parse()

	// establish a connection to the database
	client := redis.NewClient(&redis.Options{
		Addr: *redishostPtr + ":6379",
		DB:   *redisdbPtr,
	})
	_, connectErr := client.Ping().Result()
	if connectErr != nil {
		panic(connectErr)
	}
	// get the currnet number of keys to start
	keysCount := getKeysCount(client)
	// this is the main loop we are just going to continue till all the threads
	// generate the max number of files
	for keysCount < *filecountPtr {
		// if this is the first iterations we need to start by creating a file
		if keysCount == 0 {
			createFile(pathPtr, ttlPtr, client)
			// otherwise we move on to the random actions
		} else {
			// create a new random seed
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			// grap a random number to decide what we are doing this loop
			switch mode := r1.Intn(4); mode {
			case 0:
				createFile(pathPtr, ttlPtr, client)
			case 1:
				updateFile(client)
			case 2:
				readFile(client)
			case 3:
				delFile(client)
			}
		}
		// get the new count to update the control
		keysCount = getKeysCount(client)
	}
	// when we hit the limit we want to start cleaning up after ourselves
	delAllFiles(client)
}
