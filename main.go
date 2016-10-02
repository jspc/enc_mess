package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/boltdb/bolt"
    "github.com/gorilla/context"
    "golang.org/x/crypto/ssh/terminal"
)

var password []byte
var amqpUri string
var myName string
var privateKeyPath string

var storage Storage

func init() {
    flag.StringVar(&amqpUri, "amqp", "amqp://guest:guest@localhost:5671/", "URI to pass messages via")
    flag.StringVar(&myName, "sender-name", "jspc", "Simple cosmetic placeholder for sender's name")
    flag.StringVar(&privateKeyPath, "private-key", "./key.pem", "RSA Secret key PEM file")
    flag.StringVar(&storage.Path, "db-file", "./keys.db", "DB file, storing public keys. Will be created if it doesn't exist")

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

    log.Println("Starting RESTful API")

    http.HandleFunc("/", Router)
    http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))

//    return
}
