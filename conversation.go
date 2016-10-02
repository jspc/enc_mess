package main

import (
    "crypto/rsa"
    "encoding/base64"
    "log"
)

type Conversation struct {
    Sender string
    Recipient string
    EncryptionKey *rsa.PublicKey
}

func (c *Conversation)Start(sender, recipient string) {
    k := LoadPublicKey(recipient)

    c.Recipient = recipient
    c.Sender = sender
    c.EncryptionKey = k

    go func() {
        _, err := AMQPConsumer(amqpUri, recipient, sender, c)
        if err != nil {
            log.Fatal(err.Error())
        }
    }()
    select {}
}

func (c *Conversation)SendMessage(body []byte) {
    log.Println(string(body))

    msg := encrypt(body, c.EncryptionKey)
    encodedMsg := base64.StdEncoding.EncodeToString(msg)

    if err := AMQPPublisher(amqpUri, c.Recipient, encodedMsg); err != nil {
        log.Printf("Sending message to %s failed: %s", c.Recipient, err.Error())
    }

}

func (c *Conversation)ReceiveMessage(body []byte) {
    decodedMsg, err := base64.StdEncoding.DecodeString(string(body))
    if err != nil {
        log.Printf("Error decoding message from %s: %s", c.Recipient, err.Error())
    } else {
        msg := decrypt(decodedMsg)
        log.Printf("%s received '%s' from %s", c.Sender, string(msg), c.Recipient)
    }
}
