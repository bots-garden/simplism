package server

import (
	"bytes"
	"fmt"
	"log"
	simplismTypes "simplism/types"

	bolt "go.etcd.io/bbolt"
)

func initializeStoreDB(wasmArgs simplismTypes.WasmArguments, dbPath string) (*bolt.DB, error) {
	//! we can use dbPath in the future to store the data in a different location
	//👀 https://github.com/etcd-io/bbolt?tab=readme-ov-file#database-backups
	//? or find a way to save the data to a S3 bucket (recurrent backup)
	// TODO: load testing and decide if we need to add a timeout
	//db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if dbPath == "" {
		dbPath = wasmArgs.FilePath + ".store.db"
	}

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("simplism-store-bucket"))
		if err != nil {
			return fmt.Errorf("😡 When creating the store bucket: %s", err)
		}
		return nil
	})
	return db, err
}

func saveToStore(db *bolt.DB, key, value string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-store-bucket"))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

func getFromStore(db *bolt.DB, key string) string { // 🤔 should I return an error?

	var strReturnValue string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-store-bucket"))
		bufferValue := b.Get([]byte(key))

		if bufferValue == nil {
			strReturnValue = ""
		} else {
			strReturnValue = string(bufferValue)
		}
		return nil
	})
	return strReturnValue
}

func getAllFromStore(db *bolt.DB) map[string]string {
	records := map[string]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-store-bucket"))
		c := b.Cursor()

		for bufferKey, bufferValue := c.First(); bufferKey != nil; bufferKey, bufferValue = c.Next() {
			records[string(bufferKey)] = string(bufferValue)
		}

		return nil
	})
	return records
}

func getAllWitPrefixFromStore(db *bolt.DB, prefix string) map[string]string {
	records := map[string]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-store-bucket"))
		c := b.Cursor()
		bufferPrefix := []byte(prefix)

		for bufferKey, bufferValue := c.Seek(bufferPrefix); bufferKey != nil && bytes.HasPrefix(bufferKey, bufferPrefix); bufferKey, bufferValue = c.Next() {
			records[string(bufferKey)] = string(bufferValue)
		}

		return nil
	})
	return records
}

func deleteFromStore(db *bolt.DB, key string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("simplism-store-bucket"))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}

// range ?
