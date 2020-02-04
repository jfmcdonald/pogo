package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// Package file containing all the direct file IO functions used in pogo.

// Description: creates a new blank file and logs it into the redis db
// Expects: a pointers to the root path, key ttl and open redis client
// Returns: nothing
// TODO:
func createFile(pathPtr *string, ttlPtr *int, client *redis.Client) {
	// generate a random file name
	fileName := randString(10)
	// build up the full path
	filePath := *pathPtr + "/" + fileName
	// Create the file on disk
	emptyFile, err := os.Create(filePath)
	if err != nil {
		log.Println("error creating empty file at: " + filePath)
		log.Println(err)
	}
	emptyFile.Close()
	// add the place holder to the redis db
	setErr := client.Set(fileName, filePath, time.Duration(*ttlPtr)*time.Second).Err()
	if setErr != nil {
		// due to retries this should only trigger after 5 attempts
		log.Println("error setting db key for file: " + fileName)
		log.Println(setErr)
		// if we can't write the key we should delete the file to not leave orphans
		delErr := os.Remove(filePath)
		if delErr != nil {
			// if this errors too we are probably boraked
			panic("backout file deletiong failed: functon createFile")
		}
		panic("failed to create redis entry for new file")
	}
}

// Description: updates an random existing file with 8 bits of data
// Expects: pointer to an open redis client
// Returns: nothing
// TODO:
func updateFile(client *redis.Client) {
	_, val := getPath(client)
	// open our file for writing, no read we don't want to create an
	// excuse for the filesystem to cache data
	file, ioErr := os.OpenFile(val, os.O_WRONLY|os.O_APPEND, 0644)
	if ioErr != nil {
		log.Println("unable to open file for writing: " + val)
		log.Println(ioErr)
	}
	// defer the close till the end of the funciton
	defer file.Close()
	// write out an random 8 bits. We are igoring the length here since
	// we don't use it
	_, writeErr := file.WriteString(randString(8))
	if writeErr != nil {
		log.Println("unable to write out to file: " + val)
		log.Println(writeErr)
	}
}

// Description: read a file by line, generate a line count. This is just to make
// sure the file is actually read from disk in its entirety
// Expects: pointer to open redis db connection
// Returns: nothing
// TODO:
func readFile(client *redis.Client) {
	_, val := getPath(client)
	// Read in our file
	data, readErr := ioutil.ReadFile(val)
	if readErr != nil {
		log.Println("unalbe to read from file: " + val)
		log.Println(readErr)
	}
	// do something with the data so it's really in mem
	data = data
}

// Description: Delete a random file and remove it's db entery
// Expects: pointer to open redis client
// Returns: nothing
// TODO: redis key delete error check is broken
func delFile(client *redis.Client) {
	key, val := getPath(client)
	// actually remove the file.
	remErr := os.Remove(val)
	if remErr != nil {
		log.Println("unalbe to remove key: ")
		log.Println(remErr)
	}
	// remove the associated key
	client.Del(key)
	// This always errors even when it works
	//if redisDelErr != nil {
	//	panic(redisDelErr)
	//}
}

// Description: delete all remaining keys from the current run
// Expects: pointer to open redis client
// Returns: nothing
func delAllFiles(client *redis.Client) {
	// check to make sure there really are files to delete
	keysCount := getKeysCount(client)
	for keysCount > 0 {
		// why do this all twice, just call an existing function
		delFile(client)
		// update our count
		keysCount = getKeysCount(client)
	}
}
