# boltdbstore

Wrap the BoltDB. Provide methods for storing data as JSON and loading into collections.

The main goal is publish how to implement general method for unmarshal JSON into typed collection
without using the [reflect](https://golang.org/pkg/reflect/) package.


## Adding into project

To use in project import the library and set environment variable "BOLTDB_PATH".

```go
func init() {
	if p := os.Getenv("BOLTDB_PATH"); p == "" {
 		os.Setenv("BOLTDB_PATH", "test.db")
 	}
}
```


## Loading stored collection

The Stored interface defines:
 * `Bucket() []byte` should return the collection name.
 * `Next([]byte) interface{}` should implement an iterator that takes key and returns pointer to an item instance in collection.

With Next() implementation it's possible to load data into different collections.
For example let's implement support for map and slice:

```go
import "github.com/satori/go.uuid"

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

// GetRecordsMap returns collection from boltdb Bucket
func GetRecordsMap() (items RecordsMap, err error) {
	err = boltdbstore.GetStored(&items)
	return
}

// GetRecordsSlice returns collection from boltdb Bucket
func GetRecordsSlice() (items RecordsSlice, err error) {
	err = boltdbstore.GetStored(&items)
	return
}
```

For full code example take a look into tests.
