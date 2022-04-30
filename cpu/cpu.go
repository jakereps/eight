package cpu

import "fmt"

type Controller struct {
	registers []uint8
	i         uint16
	delay     uint8
	sound     uint8
	pc        uint16
	sp        uint8
	stack     []uint16
}

func (c *Controller) Vx(x uint8) uint8 {
	return c.registers[x]
}

func (c *Controller) State() string {
	return fmt.Sprintf("%+v", c)
}

func NewController() *Controller {

	return &Controller{
		registers: make([]uint8, 16),
		stack:     make([]uint16, 16),
	}
}
