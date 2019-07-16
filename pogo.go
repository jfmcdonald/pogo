package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	// command line options
	pathPtr := flag.String("path", "/tmp/pogo", "Path wher run time files will be generated")
	filecountPtr := flag.Uint64("count", 10, "Total number of files to generate")
	ttlPtr := flag.Int("ttl", 60, "Index Key/Value store default key TTL")
	redishostPtr := flag.String("dbhost", "localhost", "Hostname of the network redis server")
	redisdbPtr := flag.Int("db", 0, "redis db id you want to store keys in")
	logfilePtr := flag.String("logfile", "/dev/null", "location where you want to log message")
	redispassPtr := flag.String("pass", "", "Password for redis db")

	flag.Parse()
	// setup logging
	logFile, logErr := os.OpenFile(*logfilePtr, os.O_WRONLY|os.O_APPEND, 666)
	if logErr != nil {
		panic(logErr)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// establish a connection to the database
	client := redis.NewClient(&redis.Options{
		Addr:       *redishostPtr + ":6379",
		DB:         *redisdbPtr,
		Password:   *redispassPtr,
		MaxRetries: 5,
	})
	_, connectErr := client.Ping().Result()
	if connectErr != nil {
		log.Println("Error connecting to Redis: ")
		log.Println(connectErr)
		// we still want to panic here, we can't work with out the db
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
