package display

import "fmt"

const width = 64
const height = 32

type Text struct {
	disp [][]uint8
}

func (t *Text) State() string {
	return fmt.Sprintf("%+v", t)
}

func NewText() *Text {
	disp := make([][]uint8, height)

	for i := range disp {
		disp[i] = make([]uint8, width)
	}

	return &Text{
		disp: disp,
	}
}
