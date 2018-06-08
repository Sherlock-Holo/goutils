package websocket

import (
    "io/ioutil"
    "net/http"
    "testing"

    "fmt"
    "github.com/gorilla/websocket"
    "net/url"
    "time"
    "io"
)

var upgrader = websocket.Upgrader{}

func TestWrapper(t *testing.T) {
    done := make(chan struct{})
    go server(t, done)
    time.Sleep(time.Second)
    client(t, done)
}

func server(t *testing.T, done <-chan struct{}) {
    handle := func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)

        if err != nil {
            t.Error(err)
            return
        }

        wrapper := NewWrapper(conn)

        /*if _, err := wrapper.Write([]byte("holo")); err != nil {
            t.Error(err)
            wrapper.Close()
            return
        }*/
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

    go t.Fatal(http.ListenAndServe("127.0.0.1:9876", nil))

    <-done
}

func client(t *testing.T, done chan<- struct{}) {
    u := url.URL{Scheme: "ws", Host: "127.0.0.1:9876", Path: "/echo"}

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        t.Fatal(err)
    }

    wrapper := NewWrapper(conn)

    /*if _, err := wrapper.Write([]byte("sherlock")); err != nil {
        t.Error(err)
        wrapper.Close()
        return
    }*/
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

    close(done)
}
