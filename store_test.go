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
	"os"
	"testing"
)

func TestOpenBucket(t *testing.T) {
	db, err := openBucket([]byte(RecordsBucket))
	defer db.Close()

	if err != nil {
		t.Fatal("#TestOpenBucket")
	}
}

func TestPut(t *testing.T) {
	r := NewRecord()
	err := r.Save()
	if err != nil {
		t.Error("#TestPut,#Save", err, r)
	}

	rmap, err := GetRecordsMap()
	if err != nil {
		t.Error("#TestPut,#GetRecordsMap", err, rmap)
	}

	t.Log("#TestPut,#GetRecordsMap", rmap)
}

func TestDelete(t *testing.T) {
	r := NewRecord()
	err := r.Save()
	if err != nil {
		t.Fatal("#TestDelete,#Save", err, r)
	}

	err = DeleteRecord(r.ID)
	if err != nil {
		t.Error("#TestDelete,#DeleteRecord", err)
	}

	buf, err := GetRecord(r.ID)
	if err == nil {
		t.Error("#TestDelete,#GetRecord", err, buf)
	}
}

func TestGet(t *testing.T) {
	r := NewRecord()
	err := r.Save()
	if err != nil {
		t.Fatal("#TestGet,#Save", err, r)
	}

	buf, err := GetRecord(r.ID)
	if err != nil {
		t.Error("#TestGet,#GetRecord", err, buf)
	}

	t.Log("#TestGet,#GetRecord", buf)
}

func TestGetStored(t *testing.T) {
	r := NewRecord()
	err := r.Save()
	if err != nil {
		t.Fatal("#TestGetStored,#Save", err, r)
	}

	rslice, err := GetRecordsSlice()
	if err != nil {
		t.Error("#TestGetStored,#GetRecordsSlice", err, rslice)
	}

	rmap, err := GetRecordsMap()
	if err != nil {
		t.Error("#TestGetStored,#GetRecordsMap", err, rmap)
	}

	t.Log("#TestGetStored,#GetRecordsMap", rmap)
}

func TestMain(m *testing.M) {
	// setup
	if p := os.Getenv("BOLTDB_PATH"); p == "" {
		os.Setenv("BOLTDB_PATH", "test.db")
	}
	os.Remove(os.Getenv("BOLTDB_PATH"))
	// run tests
	e := m.Run()
	// teardown
	os.Remove(os.Getenv("BOLTDB_PATH"))
	// fin
	os.Exit(e)
}
