package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io"
    "strings"
    "time"

    "github.com/gorilla/sessions"
    "github.com/satori/go.uuid"
)

type SimpleResponse struct {
    Body interface{}
    Status int
}

type SimpleRequestBody struct {
    Body string
}

var resp SimpleResponse
var store = sessions.NewCookieStore([]byte( uuid.NewV4().String() ))
var sessionStoreID = fmt.Sprintf("enc-mess-%s", uuid.NewV4().String() )

func Router(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)

    session, err := store.Get(r, sessionStoreID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    session.Values["last-request"] = fmt.Sprintf("%s", time.Now())
    session.Values["sender"] = myName

    w = Headers(w)
    resp.Status = http.StatusOK

    conversationsPrefix := "/conversations/"

    switch  {
    case r.Method == "OPTIONS":
        resp.Body = ""


    case r.URL.Path ==  "/conversations" && r.Method == "GET":
        resp.Body = GetAllConversations(session)

    case r.URL.Path ==  "/conversations" && r.Method == "POST":
        if err, cID := AddConversation(session, r) ; err != nil {
            resp.Status = http.StatusInternalServerError
            resp.Body = err.Error()
        } else {
            session.Values["conversations"] = UpdateConversationsList(cID, session)
            resp.Body = cID
        }


    case strings.HasPrefix(r.URL.Path, conversationsPrefix) && r.Method == "GET":
        resp.Body = GetConversation(session, strings.TrimPrefix(r.URL.Path, conversationsPrefix))

    case strings.HasPrefix(r.URL.Path, conversationsPrefix) && r.Method == "POST":
        var b SimpleRequestBody
        if err := json.Unmarshal(rcToByteSlice(r.Body), &b); err != nil {
            resp.Status = http.StatusInternalServerError
            resp.Body = err.Error()
        } else {
            PostToConversation(session, strings.TrimPrefix(r.URL.Path, conversationsPrefix), []byte(b.Body) )
            resp.Body = "added"
        }


    default:
        resp.Status = http.StatusNotFound
        resp.Body = "Not found"
    }

    session.Save(r,w)
    resp.respond(w)
}

func Headers(w http.ResponseWriter) http.ResponseWriter{
    w.Header().Set("Access-Control-Allow-Headers", "requested-with, Content-Type, origin, authorization, accept, client-security-token, cache-control, x-api-key")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Max-Age", "10000")
    w.Header().Set("Cache-Control", "no-cache")

    w.Header().Set("Content-Type", "application/json")
    return w
}

func (r *SimpleResponse) respond (w http.ResponseWriter) {
    w.WriteHeader(r.Status)
    j,_ := json.Marshal(r)
    fmt.Fprintf(w, string(j))
}

func LogRequest(r *http.Request) {
    log.Printf( "%s :: %s %s",
        r.RemoteAddr,
        r.Method,
        r.URL.Path)
}

func rcToByteSlice(rc io.ReadCloser) []byte{
    buf := new(bytes.Buffer)
    buf.ReadFrom(rc)
    return buf.Bytes()
}
