package cpu

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Controller struct {
	registers []uint8
	i         uint16
	pc        uint16
	sp        uint8
	stack     *stack

	delay uint8
	dmu   *sync.RWMutex
	sound uint8
	smu   *sync.RWMutex
}

func (c *Controller) WriteVx(x uint8, b uint8) {
	c.registers[x] = b
}

func (c *Controller) Vx(x uint8) uint8 {
	return c.registers[x]
}

type ReadWriter interface {
	Read(uint16) uint8
	Write(uint16, uint8)
	Sprite(uint8) uint16
}

type Drawer interface {
	Draw(x, y uint8, data []uint8) bool
	Clear()
}

type Inputer interface {
	Pressed(uint8) bool
}

func (c *Controller) Op(r ReadWriter, d Drawer, i Inputer) {
	var stepped bool
	defer func() {
		if !stepped {
			c.pc += 2
		}
	}()

	hi := r.Read(c.pc)
	kk := r.Read(c.pc + 1)

	inst := (uint16(hi) << 8) | uint16(kk)
	nnn := inst & 0x0FFF
	n := inst & 0x000F
	x := hi & 0x0F
	y := (kk & 0xF0) >> 4

	log.Printf("instruction: %04x - hi: %02x, lo (kk): %02x, nnn: %03x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y)

	switch hi >> 4 {
	case 0x0:
		switch kk & 0xF {
		case 0x0:
			d.Clear()
		case 0xE:
			c.pc = c.stack.pop()
			c.sp -= 1
		default:
			panic(fmt.Sprintf("unknown instruction: %04x - hi: %02x, lo (kk): %02x, nnn: %03x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
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
	case 0xc:
		c.WriteVx(x, kk&uint8(rand.Intn(256)))
	case 0xd:
		data := make([]byte, 0, n)
		for i := uint16(0); i < n; i++ {
			data = append(data, r.Read(c.i+i))
		}
		collision := d.Draw(c.Vx(x), c.Vx(y), data)
		if collision {
			c.WriteVx(0xF, 0x1)
		} else {
			c.WriteVx(0xF, 0x0)
		}
	case 0xe:
		switch kk {
		case 0xa1:
			if !i.Pressed(c.Vx(x)) {
				c.pc += 2
			}
		default:
			panic(fmt.Sprintf("unknown instruction: %04x - hi: %02x, lo (kk): %02x, nnn: %03x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
		}
	case 0xf:
		switch kk {
		case 0x07:
			c.dmu.RLock()
			d := c.delay
			c.dmu.RUnlock()
			c.WriteVx(x, d)
		case 0x15:
			c.dmu.Lock()
			c.delay = c.Vx(x)
			c.dmu.Unlock()
		case 0x29:
			c.i = r.Sprite(c.Vx(x))
		case 0x33:
			data := strconv.Itoa(int(c.Vx(x)))
			for i, digit := range data {
				v, err := strconv.Atoi(string(digit))
				if err != nil {
					panic(fmt.Sprintf("failed converting to int: %s", err))
				}
				r.Write(c.i+uint16(len(data)-i), uint8(v))
			}
		case 0x65:
			var vx uint8
			for vx = 0; vx <= x; vx++ {
				log.Printf("writing - vx: %d, pos: %04x, data: %b", vx, c.i+uint16(vx), r.Read(c.i+uint16(vx)))
				c.WriteVx(vx, r.Read(c.i+uint16(vx)))
			}
		default:
			panic(fmt.Sprintf("unknown instruction: %04x - hi: %02x, lo (kk): %02x, nnn: %03x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
		}
	default:
		panic(fmt.Sprintf("unknown instruction: %04x - hi: %02x, lo (kk): %02x, nnn: %03x, n: %x, x: %d, y: %d", inst, hi, kk, nnn, n, x, y))
	}
}

func (c *Controller) State() string {
	return fmt.Sprintf("%+v", c)
}

func multiplex(n int, c <-chan time.Time) []<-chan time.Time {
	chans := make([]chan time.Time, 0, n)
	out := make([]<-chan time.Time, 0, n)
	for i := 0; i < n; i++ {
		channel := make(chan time.Time)
		chans = append(chans, channel)
		out = append(out, channel)
	}

	go func() {
		for t := range c {
			for _, cn := range chans {
				select {
				case cn <- t:
				default:
				}
			}
		}
	}()
	return out
}

func (c *Controller) startDelay(ch <-chan time.Time) {
	for range ch {
		c.dmu.Lock()
		if c.delay > 0 {
			c.delay -= 1
		}
		c.dmu.Unlock()
	}
}

func (c *Controller) startSound(ch <-chan time.Time) {
	for range ch {
		c.smu.Lock()
		if c.sound > 0 {
			c.sound -= 1
		}
		c.smu.Unlock()
	}
}

func NewController() *Controller {
	c := Controller{
		pc:        0x200,
		registers: make([]uint8, 16),
		stack:     &stack{},

		dmu: &sync.RWMutex{},
		smu: &sync.RWMutex{},
	}
	t := time.NewTicker(time.Second / 60)
	cs := multiplex(2, t.C)
	go c.startDelay(cs[0])
	go c.startSound(cs[1])
	return &c
}
