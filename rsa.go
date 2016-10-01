package main

import (
    "crypto/md5"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
    "log"
)

var block *pem.Block
var err error
var pemData,label []byte
var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func ConfigureRSA() {
    configurePrivateKey()
    configurePublicKey()
}

func configurePublicKey() {
    if pemData, err = ioutil.ReadFile(publicKeyPath); err != nil {
        log.Fatalf("Error reading pem file: %s", err)
    }

    block, _ = pem.Decode(pemData)
    if block == nil {
        log.Fatalf("failed to parse PEM block containing the public key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        log.Fatalf("failed to parse DER encoded public key: " + err.Error())
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
        publicKey = pub
    default:
        log.Fatalf("unknown type of public key")
    }
}

func configurePrivateKey() {
    if pemData, err = ioutil.ReadFile(privateKeyPath); err != nil {
        log.Fatalf("Error reading pem file: %s", err)
    }

    if block, _ = pem.Decode(pemData); block == nil || block.Type != "RSA PRIVATE KEY" {
        log.Fatal("No valid PEM data found")
    }

    s, _ := x509.DecryptPEMBlock(block, password)

    if privateKey, err = x509.ParsePKCS1PrivateKey(s); err != nil {
        log.Fatalf("Private key can't be decoded: %s", err)
    }

}

func encrypt(plaintext []byte) (encrypted []byte) {
    md5Hash := md5.New()

    if encrypted, err = rsa.EncryptOAEP(md5Hash, rand.Reader, publicKey, plaintext, label); err != nil {
        log.Fatal(err)
    }
    return
}

func decrypt(ciphertext []byte) (decrypted []byte) {
    md5Hash := md5.New()

    if decrypted, err = rsa.DecryptOAEP(md5Hash, rand.Reader, privateKey, ciphertext, label); err != nil {
        log.Fatal(err)
    }
    return
}
