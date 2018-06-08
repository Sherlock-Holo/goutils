package websocket

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "testing"
    "time"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TestWrapper(t *testing.T) {
    done := make(chan struct{})
    go func() {
        time.Sleep(time.Second)
        client(t, done)
    }()
    server(t, done)
}

func server(t *testing.T, done chan struct{}) {
    handle := func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)

        if err != nil {
            t.Error(err)
            return
        }

        wrapper := NewWrapper(conn)

        if _, err := io.WriteString(wrapper, "holo"); err != nil {
            t.Error(err)
            wrapper.Close()
            return
        }

        wrapper.CloseWrite()

        b, err := ioutil.ReadAll(wrapper)
        if err != nil {
            t.Error(err)
            wrapper.Close()
            return
        }

        fmt.Println(string(b))
        wrapper.Close()
    }

    http.HandleFunc("/echo", handle)

    go func() {
        t.Fatal(http.ListenAndServe("127.0.0.1:9876", nil))
    }()

    <-done
}

func client(t *testing.T, done chan<- struct{}) {
    defer close(done)

    u := url.URL{Scheme: "ws", Host: "127.0.0.1:9876", Path: "/echo"}

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        t.Fatal(err)
    }

    wrapper := NewWrapper(conn)

    if _, err := io.WriteString(wrapper, "sherlock"); err != nil {
        t.Error(err)
        wrapper.Close()
        return
    }

    wrapper.CloseWrite()

    b, err := ioutil.ReadAll(wrapper)
    if err != nil {
        t.Error(err)
        wrapper.Close()
        return
    }

    fmt.Println(string(b))
    wrapper.Close()
}
