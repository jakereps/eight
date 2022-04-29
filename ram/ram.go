package ram

type Storage struct {
	mem []byte
}

func NewStorage() *Storage {
	mem := make([]byte, 4096)
	copy(mem[0x050:0x0A1], presets)
	return &Storage{
		mem: mem,
	}
}
