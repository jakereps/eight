package chip8

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/jakereps/eight/cpu"
	"github.com/jakereps/eight/ram"
)

type Displayer interface {
	Draw(x, y uint8, data []byte) bool
	Clear()

	// State is used to dump the current display state while in debug mode.
	State() string
}

type Inputer interface {
	Pressed(uint8) bool
}

type Emulator struct {
	cpu   *cpu.Controller
	ram   *ram.Storage
	disp  Displayer
	keys  Inputer
	debug bool
	dstr  strings.Builder
}

func (e *Emulator) Load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for i, d := range b {
		e.ram.LoadROM(uint16(i), d)
	}
	return nil
}

func (e *Emulator) Run(ctx context.Context) error {
	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			err := e.step()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Emulator) opt() {
	e.cpu.Op(e.ram, e.disp, e.keys)
}

func (e *Emulator) step() error {
	e.opt()
	if e.debug {
		defer e.dstr.Reset()
		e.dstr.WriteString("STATE: \n")
		e.dstr.WriteString("cpu: \n")
		e.dstr.WriteString(e.cpu.State())
		e.dstr.WriteRune('\n')
		e.dstr.WriteString("ram: \n")
		e.dstr.WriteString(e.ram.State())
		e.dstr.WriteRune('\n')
		e.dstr.WriteString("disp: \n")
		e.dstr.WriteString(e.disp.State())
		e.dstr.WriteRune('\n')

		log.Println(e.dstr.String())
	}
	return nil
}

func NewEmulator(opts ...EmulatorOption) (*Emulator, error) {
	var emu Emulator
	for _, opt := range opts {
		err := opt(&emu)
		if err != nil {
			return nil, err
		}
	}
	return &emu, nil
}
