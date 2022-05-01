package inputs

import (
	"bufio"
	"os"
)

type Stdin struct {
	listener chan uint8
}

func (s *Stdin) Pressed(k uint8) bool {
	select {
	case v := <-s.listener:
		return v == k
	default:
		// nop
	}
	return false
}

func (s *Stdin) listen() {
	buf := bufio.NewReader(os.Stdin)

	for {
		r, _, err := buf.ReadRune()
		if err != nil {
			panic("invalid input")
		}
		if n := uint8(r); r&0xFFF0 == 0 {
			select {
			case s.listener <- n:
			default:
				// don't block if pressed isn't listening, just throw it away
			}
		}
	}

}

func NewStdin() *Stdin {
	listener := make(chan uint8)
	s := Stdin{
		listener: listener,
	}
	go s.listen()
	return &s
}
