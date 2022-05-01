package inputs

import (
	"bufio"
	"log"
	"os"
	"time"
)

var buf = bufio.NewReader(os.Stdin)

type Stdin struct {
	listener chan uint8
}

func (s *Stdin) Pressed(k uint8) bool {
	select {
	case v := <-s.listener:
		log.Println(v, k, mapping[k])
		return v == mapping[k]
	case <-time.After(1000 * time.Millisecond):
	}
	return false
}

func (s *Stdin) listen() {
	for {
		b, _ := buf.ReadByte()
		select {
		case s.listener <- b:
		case <-time.After(1000 * time.Millisecond):
		}
	}

}

func NewStdin() *Stdin {
	s := Stdin{
		listener: make(chan uint8, 1),
	}
	go s.listen()
	return &s
}
