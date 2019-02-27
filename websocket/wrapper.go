package websocket

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Wrapper struct {
	r          io.Reader
	conn       *websocket.Conn
	writeMutex sync.Mutex
}

func (w *Wrapper) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *Wrapper) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *Wrapper) SetDeadline(t time.Time) error {
	if err := w.SetReadDeadline(t); err != nil {
		return err
	}
	return w.SetWriteDeadline(t)
}

func (w *Wrapper) SetReadDeadline(t time.Time) error {
	return errors.WithStack(w.conn.SetReadDeadline(t))
}

func (w *Wrapper) SetWriteDeadline(t time.Time) error {
	return errors.WithStack(w.conn.SetWriteDeadline(t))
}

func NewWrapper(conn *websocket.Conn) *Wrapper {
	return &Wrapper{
		conn: conn,
	}
}

func (w *Wrapper) Write(p []byte) (n int, err error) {
	w.writeMutex.Lock()
	defer w.writeMutex.Unlock()

	if err = w.conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return
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
