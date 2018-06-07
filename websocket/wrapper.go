package websocket

import (
    "bytes"
    "io"

    "github.com/gorilla/websocket"
)

type Wrapper struct {
    conn       *websocket.Conn
    buf        bytes.Buffer
    readClosed bool
}

func NewWrapper(conn *websocket.Conn) *Wrapper {
    return &Wrapper{
        conn: conn,
    }
}

func (w *Wrapper) Read(p []byte) (n int, err error) {
    if w.buf.Len() > 0 {
        return w.buf.Read(p)
    }

    if w.readClosed {
        return 0, io.EOF
    }

    _, b, err := w.conn.ReadMessage()
    if err != nil {
        if closeError, ok := err.(*websocket.CloseError); ok {
            if closeError.Code == websocket.CloseNormalClosure {
                w.readClosed = true
                return 0, io.EOF
            }
        }

        return 0, err
    }

    w.buf.Write(b)
    return w.buf.Read(p)
}

func (w *Wrapper) Write(p []byte) (n int, err error) {
    err = w.conn.WriteMessage(websocket.BinaryMessage, p)
    if err != nil {
        return 0, err
    }

    return len(p), nil
}

func (w *Wrapper) Close() error {
    if err := w.CloseWrite(); err != nil {
        return err
    }
    return w.conn.Close()
}

func (w *Wrapper) CloseWrite() error {
    return w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
