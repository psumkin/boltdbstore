// Copyright 2016 The NorthShore Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boltdbstore

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/boltdb/bolt"
)

// Stored defines collection interface
type Stored interface {
	// Bucket returns the bucket name
	Bucket() []byte
	// Next is an iterator, takes key and returns pointer to item instance
	Next([]byte) interface{}
}

// GetStored loads all items from boltdb Bucket
func GetStored(items Stored) (err error) {
	bucket := items.Bucket()
	db, err := openBucket(bucket)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k, v []byte) error {

			if err := json.Unmarshal(v, items.Next(k)); err != nil {
				return err
			}

			return nil
		})

		return nil
	})
	return
}

// Delete deletes key from boltdb Bucket
func Delete(bucket []byte, key []byte) error {
	db, err := openBucket(bucket)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete(key)
	}); err != nil {
		return err
	}
	return nil
}

// Get gets item from boltdb Bucket
func Get(bucket []byte, key []byte, v interface{}) (err error) {
	db, err := openBucket(bucket)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		buf := b.Get(key)
		if buf == nil {
			return errors.New("Key does not exist or key is a nested bucket")
		}

		if err := json.Unmarshal(buf, &v); err != nil {
			return err
		}
		return nil
	})
	return
}

// Put puts item into boltdb Bucket as JSON
func Put(bucket []byte, key []byte, v interface{}) (err error) {
	db, err := openBucket(bucket)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		vEncoded, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := b.Put(key, vEncoded); err != nil {
			return err
		}
		return nil
	})
	return
}

func openBucket(bucket []byte) (*bolt.DB, error) {
	path := os.Getenv("BOLTDB_PATH")

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
