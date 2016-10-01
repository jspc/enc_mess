package main

import (
    "flag"
    "fmt"

    "golang.org/x/crypto/ssh/terminal"
)

var password []byte
var privateKeyPath string
var publicKeyPath string

func init() {
    flag.StringVar(&privateKeyPath, "private-key", "./key.pem", "RSA Secret key PEM file")
    flag.StringVar(&publicKeyPath, "public-key", "./key.pub", "RSA Public key PEM file")
    flag.Parse()

    fmt.Print("password: ")
    bytePassword,_ := terminal.ReadPassword(0)

    password = bytePassword
}

func main() {
    ConfigureRSA()
    plaintext := []byte("Plain text message to be encrypted")

    encrypted := encrypt(plaintext)
    decrypted := decrypt(encrypted)

    fmt.Printf("OAEP Encrypted [%s] to \n[%x]\n", string(plaintext), encrypted)
    fmt.Printf("OAEP Decrypted [%x] to \n[%s]\n", encrypted, decrypted)

    return
}
