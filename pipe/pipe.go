package pipe

import (
	"io"
)

func chanFromConn(conn io.Reader) chan []byte {
	c := make(chan []byte)

	go func() {
		b := make([]byte, 1024)

		for {
			n, err := conn.Read(b)
			if n > 0 {
				res := make([]byte, n)
				// Copy the buffer so it doesn't get changed while read by the recipient.
				copy(res, b[:n])
				c <- res
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

//Pipe to conn
func Pipe(rw1, rw2 io.ReadWriter) {
	chan1 := chanFromConn(rw1)
	chan2 := chanFromConn(rw2)

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			}
			rw2.Write(b1)
		case b2 := <-chan2:
			if b2 == nil {
				return
			}
			rw1.Write(b2)
		}
	}
}
