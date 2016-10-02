package main

import (
    "flag"
    "fmt"
    "log"
    "time"

    "github.com/boltdb/bolt"
    "golang.org/x/crypto/ssh/terminal"
)

var password []byte
var privateKeyPath string
var publicKeyPath string
var boltPath string
var amqpUri string

var storage Storage

func init() {
    flag.StringVar(&privateKeyPath, "private-key", "./key.pem", "RSA Secret key PEM file")
    flag.StringVar(&publicKeyPath, "public-key", "./key.pub", "RSA Public key PEM file")
    flag.StringVar(&storage.Path, "db-file", "./keys.db", "DB file, storing public keys. Will be created if it doesn't exist")
    flag.StringVar(&amqpUri, "amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")

    flag.Parse()

    fmt.Print("password: ")
    bytePassword,_ := terminal.ReadPassword(0)
    fmt.Println("")

    password = bytePassword
}

func main() {
    ConfigureRSA()

    boltObj, err := bolt.Open(storage.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
        log.Fatal(err)
    }

    storage.db = boltObj
    defer storage.db.Close()
    storage.Preflight()

    return
}
