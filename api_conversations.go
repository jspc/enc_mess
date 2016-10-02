package main

import (
    "encoding/json"
    "net/http"
    "strings"
    _"log"
    "github.com/gorilla/sessions"
)

type ConversationResponseItem struct {
    ID string
    Recipient string
}

type ConversationMessagesItem struct {
    ID string
    Messages []string
}

var convStore map[string]*Conversation

func GetAllConversations(s *sessions.Session) (conversations []ConversationResponseItem) {
    if str, ok := s.Values["conversations"].(string); ok {
        var c ConversationResponseItem
        for _,cid := range strings.Split(str, ",") {
            c.ID = cid
            c.Recipient = convStore[cid].Recipient

            conversations = append(conversations, c)
        }
    }
    return
}

func GetConversation(s *sessions.Session, id string) (conversation ConversationMessagesItem) {
    if str, ok := s.Values["conversations"].(string); ok {
        for _,cid := range strings.Split(str, ",") {
            conversation.ID = id
            if cid == id {
                conv := convStore[cid]
                conversation.Messages = conv.AllMessages()
            }
        }
    }
    return
}

func PostToConversation(s *sessions.Session, id string, message []byte) {
    if str, ok := s.Values["conversations"].(string); ok {
        for _,cid := range strings.Split(str, ",") {
            if cid == id {
                conv := convStore[cid]
                conv.SendMessage(message)
            }
        }
    }
    return
}

func AddConversation(s *sessions.Session, r *http.Request) (error, string) {
    var newConversation NewConversationRequest
    var id string

    if convStore == nil {
        convStore = make(map[string]*Conversation)
    }

    if err := json.Unmarshal( rcToByteSlice(r.Body), &newConversation); err != nil {
        return err, ""
    }

    if str, ok := s.Values["sender"].(string); ok {
        var c Conversation
        c.Initialise(str, newConversation.Recipient)
        go c.Start()

        id = c.ID
        convStore[id] = &c
    }
    return nil, id
}

func UpdateConversationsList(c string, s *sessions.Session) (conversations string) {
    var cList []string

    if s.Values["conversations"] == nil {
        conversations = c

    } else if str, ok := s.Values["conversations"].(string); ok {
        cList = strings.Split(str, ",")
        cList = append(cList, c)
        conversations = strings.Join(cList, ",")
    }
    return
}
