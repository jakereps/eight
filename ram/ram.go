package ram

import "fmt"

const (
	spriteStart = 0x050
	spriteEnd   = 0x0A1
)

type Storage struct {
	mem []uint8
}

func (s *Storage) Read(loc uint16) uint8 {
	return s.mem[loc]
}

func (s *Storage) Write(loc uint16, b uint8) {
	s.mem[loc] = b
}

func (s *Storage) LoadROM(loc uint16, b uint8) {
	s.Write(0x200+uint16(loc), b)
}

func (s *Storage) Sprite(rep uint8) uint16 {
	return spriteStart + (uint16(rep) * 5)

}

func (s *Storage) State() string {
	return fmt.Sprintf("%+v", s)
}

func NewStorage() *Storage {
	mem := make([]uint8, 4096)
	copy(mem[spriteStart:spriteEnd], presets)
	return &Storage{
		mem: mem,
	}
}
