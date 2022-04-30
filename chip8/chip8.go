package chip8

import (
	"log"
	"strings"

	"github.com/jakereps/eight/cpu"
	"github.com/jakereps/eight/ram"
)

type Display interface {

	// State is used to dump the current display state while in debug mode.
	State() string
}

type Emulator struct {
	cpu   *cpu.Controller
	ram   *ram.Storage
	disp  Display
	debug bool
	dstr  strings.Builder
}

func (e *Emulator) step() {

	if e.debug {
		defer e.dstr.Reset()
		e.dstr.WriteString("STATE: \n")
		e.dstr.WriteString("cpu: \n")
		e.dstr.WriteString(e.cpu.State())
		e.dstr.WriteString("ram: \n")
		e.dstr.WriteString(e.ram.State())
		e.dstr.WriteString("disp: \n")
		e.dstr.WriteString(e.disp.State())
		e.dstr.WriteRune('\n')

		log.Printf(e.dstr.String())
	}
}

func (e *Emulator) Run() error {
	return nil
}

func NewEmulator(debug bool) *Emulator {
	emu := Emulator{}
	if debug {
		emu.dstr = strings.Builder{}
	}
	return &emu
}
