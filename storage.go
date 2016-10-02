package main

import (
    "log"
    "github.com/boltdb/bolt"
)

type Storage struct {
    db *bolt.DB
    bucket []byte
    Path string
}

var bucket *bolt.Bucket

func (s *Storage)Preflight(){
    s.bucket = bs("pub")
}

func (s *Storage)AddKey(id string, keyData string) {
    err := s.db.Update(func(tx *bolt.Tx) error {

        if bucket, err = tx.CreateBucketIfNotExists(s.bucket) ; err != nil {
            panic( err )
        }

        if err = bucket.Put(bs(id), bs(keyData)) ; err != nil {
            panic( err )
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Error adding data: %s", err.Error())
    }
}

func (s *Storage)GetKey(id string) (k []byte){
    err := s.db.View(func(tx *bolt.Tx) error {
        bucket = tx.Bucket(s.bucket)

        k = bucket.Get(bs(id))
        return nil
    })

    if err != nil {
        panic(err)
        log.Fatalf(err.Error())
    }

    return k
}

func bs(s string) []byte {
    return []byte(s)
}
