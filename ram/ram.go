package ram

import "fmt"

type Storage struct {
	mem []uint8
}

func (s *Storage) Load(loc uint16) uint8 {
	return s.mem[loc]
}

func (s *Storage) Set(loc uint16, b uint8) {
	s.mem[loc] = b
}

func (s *Storage) LoadROM(loc uint16, b uint8) {
	s.Set(0x200+uint16(loc), b)
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
