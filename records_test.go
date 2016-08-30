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
	"log"

	"github.com/satori/go.uuid"
)

// RecordsBucket defines boltdb bucket for example data
const RecordsBucket = "Records"

// Record represents example data
type Record struct {
	ID uuid.UUID
}

// RecordsMap represents example data
type RecordsMap map[string]*Record

// RecordsSlice represents example data
type RecordsSlice []*Record

// Used to avoid recursion in UnmarshalJSON below
type record Record

// NewRecord implements constructor
func NewRecord() *Record {
	return &Record{uuid.NewV4()}
}

// Save saves item in boltdb Bucket as JSON
func (r Record) Save() error {
	return Put([]byte(RecordsBucket), []byte(r.ID.String()), r)
}

// UnmarshalJSON implements custom unmarshaller for typed unmarshalling check
func (r *Record) UnmarshalJSON(b []byte) (err error) {
	log.Println("#UnmarshalJSON")
	buf := record{}
	err = json.Unmarshal(b, &buf)
	if err == nil {
		*r = Record(buf)
	}
	return
}

// Bucket implements Stored interface
func (RecordsMap) Bucket() []byte {
	return []byte(RecordsBucket)
}

// Next implements Stored interface
func (items *RecordsMap) Next(k []byte) interface{} {
	// Check for assignment to entry in nil map
	if *items == nil {
		*items = make(RecordsMap)
	}

	(*items)[string(k)] = &Record{}
	return (*items)[string(k)]
}

// Bucket implements Stored interface
func (RecordsSlice) Bucket() []byte {
	return []byte(RecordsBucket)
}

// Next implements Stored interface
func (items *RecordsSlice) Next([]byte) interface{} {
	*items = append(*items, &Record{})
	return &(*items)[len(*items)-1]
}

// DeleteRecord deletes item from boltdb Bucket
func DeleteRecord(id uuid.UUID) error {
	return Delete([]byte(RecordsBucket), []byte(id.String()))
}

// GetRecord gets item from boltdb Bucket
func GetRecord(id uuid.UUID) (r Record, err error) {
	err = Get([]byte(RecordsBucket), []byte(id.String()), &r)
	return
}

// GetRecordsMap returns collection from boltdb Bucket
func GetRecordsMap() (items RecordsMap, err error) {
	err = GetStored(&items)
	return
}

// GetRecordsSlice returns collection from boltdb Bucket
func GetRecordsSlice() (items RecordsSlice, err error) {
	err = GetStored(&items)
	return
}
