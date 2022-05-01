package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/jakereps/eight/chip8"
	"github.com/jakereps/eight/cpu"
	"github.com/jakereps/eight/display"
	"github.com/jakereps/eight/inputs"
	"github.com/jakereps/eight/ram"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("failed running emulator: %s", err)
	}
}

func run() error {
	rom := flag.String("rom", "", "the rom to load/play")
	debug := flag.Bool("debug", false, "run in debug mode (dumps state on steps)")
	flag.Parse()

	if *rom == "" {
		return errors.New("no rom specified")
	}

	opts := []chip8.EmulatorOption{
		chip8.WithDisplay(display.NewText()),
		chip8.WithInputer(inputs.NewStdin()),
		chip8.WithCPU(cpu.NewController()),
		chip8.WithRAM(ram.NewStorage()),
	}

	if *debug {
		opts = append(opts, chip8.WithDebug())
	}

	emu, err := chip8.NewEmulator(opts...)
	if err != nil {
		return err
	}

	err = emu.Load(*rom)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go emu.Run(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		cancel()
	case <-ctx.Done():
	}
	return nil
}
