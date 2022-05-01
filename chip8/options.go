package chip8

import (
	"strings"

	"github.com/jakereps/eight/cpu"
	"github.com/jakereps/eight/ram"
)

type EmulatorOption func(*Emulator) error

func WithDisplay(disp Display) EmulatorOption {
	return func(e *Emulator) error {
		e.disp = disp
		return nil
	}
}

func WithCPU(cpu *cpu.Controller) EmulatorOption {
	return func(e *Emulator) error {
		e.cpu = cpu
		return nil
	}
}

func WithRAM(ram *ram.Storage) EmulatorOption {
	return func(e *Emulator) error {
		e.ram = ram
		return nil
	}
}

func WithDebug() EmulatorOption {
	return func(e *Emulator) error {
		e.debug = true
		e.dstr = strings.Builder{}
		return nil
	}
}
