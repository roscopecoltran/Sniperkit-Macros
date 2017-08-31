package plugins

import (
	"github.com/kohkimakimoto/bucketstore"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"github.com/boltdb/bolt"
)

var fileName string

func BucketStore() {
	db, err := bucketstore.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bucket := db.Bucket("MyBucket")

	// put key/value item
	err = bucket.PutRaw([]byte("user001"), []byte(`{"name": "kohkimakimoto", "age": 36}`))
	if err != nil {
		panic(err)
	}

	// get value
	v, err := bucket.GetRaw([]byte("user001"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(v))
	// {"age":36,"name":"kohkimakimoto"}
}

func BucketStoreSearch() {
	// open database
	db, err := bucketstore.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bucket := db.Bucket("MyBucket")

	// put data (ignore errors)
	bucket.PutRaw([]byte("user001"), []byte(`{"name": "hoge", "age": 20}`))
	bucket.PutRaw([]byte("user002"), []byte(`{"name": "foo", "age": 31}`))
	bucket.PutRaw([]byte("user003"), []byte(`{"name": "bar", "age": 18}`))
	bucket.PutRaw([]byte("user004"), []byte(`{"name": "aaa", "age": 40}`))
	bucket.PutRaw([]byte("user005"), []byte(`{"name": "xxx", "age": 41}`))
	bucket.PutRaw([]byte("user006"), []byte(`{"name": "ccc", "age": 50}`))

    // query
	q := bucket.Query()
	q.Filter = &bucketstore.PropValueRangeFilter{
		Property: "age",
		Min: 20,
		Max: 40,
	}
	items, err := q.AsList()
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		fmt.Println(string(item.Key), string(item.Value))
	}
	// user001 {"age":20,"name":"hoge"}
	// user002 {"age":31,"name":"foo"}
	// user004 {"age":40,"name":"aaa"}
}

func BucketStoreOpenFile() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No db file provided")
		os.Exit(1)
	}
	fileName := args[1]
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		fmt.Println("Can't open db", err)
		os.Exit(1)
	}
	dump(db)
}

func BucketStoreDump(db *bolt.DB) {
	data := map[string]interface{}{}
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			data[string(name)] = readBucket(b)
			return nil
		})
		return nil
	})
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Can't marshal data to JSON", err)
		os.Exit(1)
	}
	dumpFileName := path.Base(fileName)
	dumpFileName = strings.TrimSuffix(dumpFileName, path.Ext(dumpFileName)) + "_dump.json"
	err = ioutil.WriteFile(dumpFileName, encoded, 0644)
	if err != nil {
		fmt.Println("Can't write dump file", err)
		os.Exit(1)
	}
	fmt.Println("Database dumped to file " + dumpFileName)
}

func BucketStoreReadBucket(b *bolt.Bucket) map[string]interface{} {
	data := map[string]interface{}{}
	b.ForEach(func(k, v []byte) error {
		if subB := b.Bucket(k); subB != nil {
			data[string(k)] = readBucket(subB)
			return nil
		}
		var _data interface{}
		err := json.Unmarshal(v, &_data)
		if err != nil {
			fmt.Println("Can't unmarshal data", string(v), err)
			return nil
		}
		data[string(k)] = _data
		return nil
	})
	return data
}
