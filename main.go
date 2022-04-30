package main

import (
	"errors"
	"flag"
	"log"

	"github.com/jakereps/eight/chip8"
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

	emu := chip8.NewEmulator(*debug)

	return emu.Run()
}
