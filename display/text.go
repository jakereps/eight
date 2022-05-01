package display

import (
	"fmt"
	"strings"
)

const width = 64
const height = 32

type Text struct {
	disp [][]uint8
}

func (t *Text) Clear() {
	t.disp = newDisp()
}

func (t *Text) Draw(x, y uint8, data []byte) bool {
	var collision bool
	for i := range data {
		for j := 0; j < 8; j++ {
			pixel := (0b1000_0000 >> j) & data[i]
			if pixel == 0 {
				continue
			}

			dx := int(x) + j
			if dx >= width {
				dx %= 64
			}
			dy := int(y) + i
			if dy >= height {
				dy %= 32
			}
			// log.Printf("drawing - x: %d, y: %d, byte: %b i: %d, j: %d", dx, y, data[i], i, j)

			last := t.disp[dy][dx]
			t.disp[dy][dx] ^= 1
			if last != 0 && t.disp[dy][dx] == 0 {
				collision = true
			}
		}
	}
	// fmt.Println(t.State())
	return collision
}

func (t *Text) State() string {
	var out strings.Builder
	for y := range t.disp {
		for x := range t.disp[y] {
			out.WriteString(fmt.Sprintf("%b", t.disp[y][x]))
		}
		out.WriteRune('\n')
	}
	return out.String()
}

func newDisp() [][]uint8 {
	disp := make([][]uint8, height)
	for i := range disp {
		disp[i] = make([]uint8, width)
	}
	return disp
}

func NewText() *Text {
	return &Text{
		disp: newDisp(),
	}
}
