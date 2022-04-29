package bus

import "github.com/jakereps/eight/cpu"

type Bus struct {
	cpu *cpu.Controller
}

func NewBus() *Bus {
	return &Bus{}
}
