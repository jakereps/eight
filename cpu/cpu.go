package cpu

import (
	"fmt"
	"log"
)

type Controller struct {
	registers []uint8
	i         uint16
	delay     uint8
	sound     uint8
	pc        uint16
	sp        uint8
	stack     *stack
}

func (c *Controller) WriteVx(x uint8, b uint8) {
	c.registers[x] = b
}

func (c *Controller) Vx(x uint8) uint8 {
	return c.registers[x]
}

type Loader interface {
	Load(uint16) uint8
}

type Drawer interface {
	Draw(x, y uint8, data []uint8) bool
	Clear()
}

func (c *Controller) Op(l Loader, d Drawer) {
	var stepped bool
	defer func() {
		if !stepped {
			c.pc += 2
		}
	}()

	hi := l.Load(c.pc)
	kk := l.Load(c.pc + 1)

	inst := (uint16(hi) << 8) | uint16(kk)
	nnn := inst & 0x0FFF
	n := inst & 0x000F
	x := hi & 0x0F
	y := (kk & 0xF0) >> 4

	log.Printf("instruction: %x - hi: %x, lo (kk): %x, nnn: %x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y)

	switch hi >> 4 {
	case 0x0:
		switch kk & 0xF {
		case 0x0:
			d.Clear()
		case 0xE:
			c.pc = c.stack.pop()
			c.sp -= 1
			stepped = true
		default:
			panic(fmt.Sprintf("unknown instruction: %x - hi: %x, lo (kk): %x, nnn: %x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
		}
	case 0x1:
		c.pc = nnn
		stepped = true
	case 0x2:
		c.sp += 1
		c.stack.push(c.pc)
		c.pc = nnn
		stepped = true
	case 0x3:
		if c.Vx(x) == kk {
			c.pc += 2
		}
	case 0x4:
		if c.Vx(x) != kk {
			c.pc += 2
		}
	case 0x5:
		if c.Vx(x) == c.Vx(y) {
			c.pc += 2
		}
	case 0x6:
		c.WriteVx(x, kk)
	case 0x7:
		c.WriteVx(x, c.Vx(x)+kk)
	case 0x9:
		if c.Vx(x) != c.Vx(y) {
			c.pc += 2
		}
	case 0xa:
		c.i = nnn
	case 0xd:
		data := make([]byte, 0, n)
		for i := uint16(0); i < n; i++ {
			data = append(data, l.Load(c.i+i))
		}
		collision := d.Draw(c.Vx(x), c.Vx(y), data)
		if collision {
			c.WriteVx(0xF, 0x1)
		} else {
			c.WriteVx(0xF, 0x0)
		}
	default:
		panic(fmt.Sprintf("unknown instruction: %x - hi: %x, lo (kk): %x, nnn: %x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
	}
}

func (c *Controller) State() string {
	return fmt.Sprintf("%+v", c)
}

func NewController() *Controller {

	return &Controller{
		pc:        0x200,
		registers: make([]uint8, 16),
		stack:     &stack{},
	}
}
