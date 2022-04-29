package chip8

type Emulator struct{}

func (e *Emulator) Run() error {
	return nil
}

func NewEmulator() *Emulator {
	return &Emulator{}
}
