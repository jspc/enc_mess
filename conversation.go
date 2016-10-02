package main

import (
    "crypto/rsa"
    "encoding/base64"
    "fmt"
    "log"
    "time"

    "github.com/satori/go.uuid"
)

type Message struct {
    Time int
    Line string
    Who string
    Failed bool
}

type Conversation struct {
    EncryptionKey *rsa.PublicKey
    ID string
    Recipient string
    Sender string
    Messages []Message
}

// As unmarshal'd in api.go
type NewConversationRequest struct {
    Recipient string
    Sender string
}

func (c *Conversation)Initialise(sender, recipient string) {
    k := LoadPublicKey(recipient)

    c.EncryptionKey = k
    c.ID = uuid.NewV4().String()
    c.Recipient = recipient
    c.Sender = sender
}

func (c *Conversation)Start() {
    _, err := AMQPConsumer(amqpUri, c)
    if err != nil {
        log.Fatal(err.Error())
    }
}

func (c *Conversation)SendMessage(body []byte) {
    msgLine := logLine(c.Sender, body)

    msg := encrypt(body, c.EncryptionKey)
    encodedMsg := base64.StdEncoding.EncodeToString(msg)

    if err := AMQPPublisher(amqpUri, c.Recipient, encodedMsg); err != nil {
        log.Printf("Sending message to %s failed: %s", c.Recipient, err.Error())
        msgLine.Failed = true
    }

    c.Messages = append(c.Messages, msgLine)
}

func (c *Conversation)ReceiveMessage(body []byte) {
    decodedMsg, err := base64.StdEncoding.DecodeString(string(body))
    if err != nil {
        log.Printf("Error decoding message from %s: %s", c.Recipient, err.Error())
    } else {
        msg := decrypt(decodedMsg)
        msgLine := logLine(c.Recipient, msg)

        c.Messages = append(c.Messages, msgLine)
    }
}

func (c *Conversation)AllMessages()(messages []string) {
    for _,m := range c.Messages {
        messages = append(messages, m.String())
    }
    return
}

func logLine(who string, msg []byte) (m Message) {
    m.Time = int(time.Now().Unix())
    m.Line = string(msg)
    m.Who = who
    m.Failed = false

    return
}

func (m *Message)String() string {
    return fmt.Sprintf("%d : %s - %s", m.Time, m.Who, m.Line)
}
