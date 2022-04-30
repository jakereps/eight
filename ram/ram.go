package ram

import "fmt"

type Storage struct {
	mem []uint8
}

func (s *Storage) Load(loc int) uint8 {
	return s.mem[loc]
}

func (s *Storage) Set(loc int, b uint8) {
	s.mem[loc] = b
}

func (s *Storage) State() string {
	return fmt.Sprintf("%+v", s)
}

func NewStorage() *Storage {
	mem := make([]uint8, 4096)
	copy(mem[0x050:0x0A1], presets)
	return &Storage{
		mem: mem,
	}
}
