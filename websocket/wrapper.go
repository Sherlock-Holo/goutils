package websocket

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

type Wrapper struct {
	r          io.Reader
	conn       *websocket.Conn
	writeMutex sync.Mutex
}

func NewWrapper(conn *websocket.Conn) *Wrapper {
	return &Wrapper{
		conn: conn,
	}
}

func (w *Wrapper) Write(p []byte) (n int, err error) {
	w.writeMutex.Lock()
	defer w.writeMutex.Unlock()

	if writer, err := w.conn.NextWriter(websocket.BinaryMessage); err != nil {
		return 0, err
	} else {
		if _, err := writer.Write(p); err != nil {
			return 0, err
		}

		if err := writer.Close(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (w *Wrapper) Read(p []byte) (n int, err error) {
	for {
		if w.r == nil {
			// Advance to next message.
			_, w.r, err = w.conn.NextReader()
			if err != nil {
				if closeError, ok := err.(*websocket.CloseError); ok {
					if closeError.Code == websocket.CloseNormalClosure {
						return 0, io.EOF
					}
					return 0, err
				}
				return 0, err
			}
		}
		n, err = w.r.Read(p)
		if err == io.EOF {
			// At end of message.
			w.r = nil
			if n > 0 {
				return n, nil
			} else {
				// No data read, continue to next message.
				continue
			}
		}
		return n, err
	}
}

func (w *Wrapper) Close() error {
	if err := w.CloseWrite(); err != nil {
		return err
	}
	return w.conn.Close()
}

func (w *Wrapper) CloseWrite() error {
	w.writeMutex.Lock()
	defer w.writeMutex.Unlock()

	return w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
